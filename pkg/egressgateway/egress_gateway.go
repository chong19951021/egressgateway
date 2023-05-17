// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package egressgateway

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/spidernet-io/egressgateway/pkg/config"
	"github.com/spidernet-io/egressgateway/pkg/constant"
	egress "github.com/spidernet-io/egressgateway/pkg/k8s/apis/egressgateway.spidernet.io/v1beta1"
	"github.com/spidernet-io/egressgateway/pkg/utils"
)

const (
	indexEgressNodeEgressGateway = "egressNodeEgressGatewayIndex"
	indexNodeEgressGateway       = "nodeEgressGatewayIndex"
)

type egnReconciler struct {
	client client.Client
	log    *zap.Logger
	config *config.Config
}

func (r egnReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	kind, newReq, err := utils.ParseKindWithReq(req)
	if err != nil {
		r.log.Sugar().Infof("parse req(%v) with error: %v", req, err)
		return reconcile.Result{}, err
	}
	log := r.log.With(
		zap.String("namespacedName", newReq.NamespacedName.String()),
		zap.String("kind", kind),
	)
	log.Info("reconciling")
	switch kind {
	case "EgressGateway":
		return r.reconcileEG(ctx, newReq, log)
	case "EgressGatewayPolicy":
		return r.reconcileEGP(ctx, newReq, log)
	case "Node":
		return r.reconcileNode(ctx, newReq, log)
	case "EgressNode":
		return r.reconcileEN(ctx, newReq, log)
	default:
		return reconcile.Result{}, nil
	}
}

// reconcileNode reconcile node
func (r egnReconciler) reconcileNode(ctx context.Context, req reconcile.Request, log *zap.Logger) (reconcile.Result, error) {
	deleted := false
	node := new(corev1.Node)
	err := r.client.Get(ctx, req.NamespacedName, node)
	if err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
		deleted = true
	}
	deleted = deleted || !node.GetDeletionTimestamp().IsZero()

	egList := &egress.EgressGatewayList{}
	if err := r.client.List(context.Background(), egList); err != nil {
		return reconcile.Result{Requeue: true}, nil
	}

	// 节点删除。节点 NoReady 事件，则在 EgressNode 事件中完成
	if deleted {
		r.log.Info("request item is deleted")
		err := r.deleteNodeFromEGs(ctx, req.Name, egList)
		if err != nil {
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Node 检查标签修改
	// 1、listEG 一一进行对比，查看node 情况是否发生变化
	// 2、新增的则添加一个空的节点信息，需要删除的，则重新分配再更新
	for _, eg := range egList.Items {
		selNode, err := metav1.LabelSelectorAsSelector(eg.Spec.NodeSelector.Selector)
		if err != nil {
			return reconcile.Result{Requeue: true}, nil
		}
		isMatch := selNode.Matches(labels.Set(node.Labels))
		if isMatch {
			// 检查，status 是否存在该node 信息。不存在则新增一个空的
			_, isExist := GetEIPStatusByNode(node.Name, eg)
			if !isExist {
				eg.Status.NodeList = append(eg.Status.NodeList, egress.EgressIPStatus{Name: node.Name})

				r.log.Sugar().Debugf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
				err := r.client.Status().Update(ctx, &eg)
				if err != nil {
					r.log.Sugar().Errorf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
					return reconcile.Result{Requeue: true}, nil
				}
			}
		} else {
			// 检查，status 中是否存在该node 信息。存在则删除，在重新分配相应的 policy
			_, isExist := GetEIPStatusByNode(node.Name, eg)
			if isExist {
				err := r.deleteNodeFromEG(ctx, node.Name, eg)
				if err != nil {
					return reconcile.Result{Requeue: true}, nil
				}
			}
		}
	}

	return reconcile.Result{}, nil
}

// reconcileEG reconcile egress gateway
func (r egnReconciler) reconcileEG(ctx context.Context, req reconcile.Request, log *zap.Logger) (reconcile.Result, error) {
	deleted := false
	eg := &egress.EgressGateway{}
	err := r.client.Get(ctx, req.NamespacedName, eg)
	if err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{Requeue: true}, err
		}
		deleted = true
	}
	deleted = deleted || !eg.GetDeletionTimestamp().IsZero()

	if deleted {
		log.Info("request item is deleted")
		return reconcile.Result{}, nil
	}

	if eg.Spec.NodeSelector.Selector == nil {
		log.Info("nodeSelector is nil, skip reconcile")
		return reconcile.Result{}, nil
	}

	// diff nodeSelector
	// 1、获取当前最新符合条件的 node
	// 2、通过 status 的nodeList 拿到旧的节点信息
	// 3、比较需要删除的节点
	// 4、待删除节点有哪些 policy 的网关节点，重新选网关节点
	// 5、EIP 再根据 policy 中的配置来生效。也就是说待删除节点上所有的 EIP 都会重新分配

	// 1、获取当前最新符合条件的 node
	newNodeList := &corev1.NodeList{}
	selNodes, err := metav1.LabelSelectorAsSelector(eg.Spec.NodeSelector.Selector)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = r.client.List(ctx, newNodeList, &client.ListOptions{
		LabelSelector: selNodes,
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Sugar().Debugf("number of selected nodes: %d", len(newNodeList.Items))

	// 2、拿到需要删除的节点
	var delNodeList []egress.EgressIPStatus
	for _, oldNode := range eg.Status.NodeList {
		isDel := true
		for _, node := range newNodeList.Items {
			if oldNode.Name == node.Name {
				isDel = false
				break
			}
		}

		if isDel {
			delNodeList = append(delNodeList, oldNode)
		}
	}
	log.Sugar().Debugf("delete a gateway nodes: %d", delNodeList)

	// 重新分配网关节点
	if len(delNodeList) != 0 {
		// perNodeListMap 保存最新的 EgressGateway.status.nodeList
		perNodeListMap := make(map[string]egress.EgressIPStatus, 0)
		// 初始化 perNodeListMap
		for _, node := range eg.Status.NodeList {
			isDel := false
			for _, item := range delNodeList {
				if item.Name == node.Name {
					isDel = true
					break
				}
			}

			if !isDel {
				perNodeListMap[node.Name] = node
			}
		}

		for _, node := range newNodeList.Items {
			_, ok := perNodeListMap[node.Name]
			if !ok {
				perNodeListMap[node.Name] = egress.EgressIPStatus{Name: node.Name}
			}
		}

		// 拿出需要重新选择网关节点的 policy
		var reSetPolicies []egress.Policy
		for _, item := range delNodeList {
			for _, eip := range item.Eips {
				reSetPolicies = append(reSetPolicies, eip.Policies...)
			}
		}

		// 逐个 policy 重新分配
		for _, policy := range reSetPolicies {
			err = r.reAllocatorPolicy(ctx, policy, eg, perNodeListMap)
			if err != nil {
				log.Sugar().Errorf("reallocator Failed to reassign a gateway node for EgressGatewayPolicy %v: %v", policy, err)
				return reconcile.Result{Requeue: true}, err
			}
		}

		var perNodeList []egress.EgressIPStatus
		for _, node := range perNodeListMap {
			perNodeList = append(perNodeList, node)
		}
		eg.Status.NodeList = perNodeList

		log.Sugar().Debugf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
		err = r.client.Status().Update(ctx, eg)
		if err != nil {
			log.Sugar().Errorf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
			return reconcile.Result{Requeue: true}, err
		}
	}

	return reconcile.Result{}, nil
}

// reconcileEG reconcile egress node
func (r egnReconciler) reconcileEN(ctx context.Context, req reconcile.Request, log *zap.Logger) (reconcile.Result, error) {
	deleted := false
	log.Sugar().Infof("reconcileNode: Delete %v event", req.Name)
	en := new(egress.EgressNode)
	en.Name = req.Name
	err := r.client.Delete(ctx, en)
	if err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{Requeue: true}, err
		}
	}
	deleted = deleted || !en.GetDeletionTimestamp().IsZero()

	// 已经处理了 node 的删除事件，所以这里不用处理
	if deleted {
		log.Info("request item is deleted")
		return reconcile.Result{}, nil
	}

	// 非success 状态的节点，将该节点上的 policy  重新分配
	if en.Status.Phase != egress.EgressNodeSucceeded {
		egList := &egress.EgressGatewayList{}
		if err := r.client.List(context.Background(), egList); err != nil {
			return reconcile.Result{Requeue: true}, nil
		}
		for _, eg := range egList.Items {
			policies, isExist := GetEIPStatusByNode(en.Name, eg)
			if isExist {
				perNodeListMap := make(map[string]egress.EgressIPStatus, 0)

				// 初始化 perNodeListMap
				for _, node := range eg.Status.NodeList {
					if node.Name == en.Name {
						continue
					} else {
						perNodeListMap[node.Name] = node
					}
				}

				for _, policy := range policies {
					err = r.reAllocatorPolicy(ctx, policy, &eg, perNodeListMap)
					if err != nil {
						log.Sugar().Errorf("reallocator Failed to reassign a gateway node for EgressGatewayPolicy %v: %v", policy, err)
						return reconcile.Result{Requeue: true}, err
					}
				}

				var perNodeList []egress.EgressIPStatus
				for _, node := range perNodeListMap {
					perNodeList = append(perNodeList, node)
				}

				perNodeList = append(perNodeList, egress.EgressIPStatus{Name: en.Name})
				eg.Status.NodeList = perNodeList

				log.Sugar().Debugf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
				err = r.client.Status().Update(ctx, &eg)
				if err != nil {
					log.Sugar().Errorf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
					return reconcile.Result{Requeue: true}, err
				}
			}
		}

	}

	return reconcile.Result{}, nil
}

// reconcileEN reconcile egress gateway policy
func (r egnReconciler) reconcileEGP(ctx context.Context, req reconcile.Request, log *zap.Logger) (reconcile.Result, error) {
	deleted := false
	isUpdete := false
	egp := &egress.EgressGatewayPolicy{}
	err := r.client.Get(ctx, req.NamespacedName, egp)
	if err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{Requeue: true}, err
		}
		deleted = true
	}

	deleted = deleted || !egp.GetDeletionTimestamp().IsZero()

	egName := egp.Spec.EgressGatewayName
	eg := &egress.EgressGateway{}
	err = r.client.Get(ctx, types.NamespacedName{Name: egName}, eg)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
		log.Sugar().Errorf("get EgressGateway err:%v", err)
		return reconcile.Result{Requeue: true}, err
	}

	policy := egress.Policy{Name: req.Name, Namespace: req.Namespace}
	if deleted {
		// 将EgressGateway 中该 policy的信息删除。如果引用的 EIP，已无其他policy 使用，则进行回收。
		DeletePlocyFromEG(policy, eg)

		log.Sugar().Debugf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
		err = r.client.Status().Update(ctx, eg)
		if err != nil {
			log.Sugar().Errorf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
			return reconcile.Result{Requeue: true}, err
		}
	}

	// policy 是否有网关节点，没有则选择
	eipStatus, isExist := GetEIPStatusByPolicy(policy, *eg)
	if !isExist {
		// 重新分配
		perNodeListMap := make(map[string]egress.EgressIPStatus, 0)
		for _, item := range eg.Status.NodeList {
			perNodeListMap[item.Name] = item
		}

		err := r.reAllocatorPolicy(ctx, policy, eg, perNodeListMap)
		if err != nil {
			r.log.Sugar().Errorf("reallocator Failed to reassign a gateway node for EgressGatewayPolicy %v: %v", policy, err)
			return reconcile.Result{Requeue: true}, err
		}

		var perNodeList []egress.EgressIPStatus
		for _, node := range perNodeListMap {
			perNodeList = append(perNodeList, node)
		}
		eg.Status.NodeList = perNodeList

		isUpdete = true
		goto update
	} else {
		// 如果指定了 EIP，如果 EIP 不对，则纠正。使用节点IP，相当于指定 EIP 为空
		for i, eip := range eipStatus.Eips {
			for j, p := range eip.Policies {
				if p == policy {
					isReAllocatorPolicy := false
					if egp.Spec.EgressIP.UseNodeIP && (eip.IPv4 != "" || eip.IPv6 != "") {
						isReAllocatorPolicy = true
					} else if egp.Spec.EgressIP.IPv4 != "" && egp.Spec.EgressIP.IPv4 != eip.IPv4 {
						isReAllocatorPolicy = true
					} else if egp.Spec.EgressIP.IPv6 != "" && egp.Spec.EgressIP.IPv6 != eip.IPv6 {
						isReAllocatorPolicy = true
					}

					if isReAllocatorPolicy {
						// 当前的 EIP 与指定的不一致。需要重新分配
						eipStatus.Eips[i].Policies = append(eipStatus.Eips[i].Policies[:j], eipStatus.Eips[i].Policies[j+1:]...)
						perNodeListMap := make(map[string]egress.EgressIPStatus, 0)
						// 初始化 perNodeListMap
						for _, node := range eg.Status.NodeList {
							if node.Name == eipStatus.Name {
								perNodeListMap[node.Name] = eipStatus
							} else {
								perNodeListMap[node.Name] = node
							}
						}

						err := r.reAllocatorPolicy(ctx, policy, eg, perNodeListMap)
						if err != nil {
							r.log.Sugar().Errorf("reallocator Failed to reassign a gateway node for EgressGatewayPolicy %v: %v", policy, err)
							return reconcile.Result{Requeue: true}, err
						}

						var perNodeList []egress.EgressIPStatus
						for _, node := range perNodeListMap {
							perNodeList = append(perNodeList, node)
						}
						eg.Status.NodeList = perNodeList
					}

					isUpdete = true
					goto update
				}
			}
		}

	}

update:
	if isUpdete {
		r.log.Sugar().Debugf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
		err = r.client.Status().Update(ctx, eg)
		if err != nil {
			r.log.Sugar().Errorf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
			return reconcile.Result{Requeue: true}, err
		}
	}

	return reconcile.Result{}, nil
}

// 将该节点从 EGList 中删除
func (r egnReconciler) deleteNodeFromEGs(ctx context.Context, nodeName string, egList *egress.EgressGatewayList) error {
	// 1、找出选择该节点作为网关节点的 EgressGateway
	// 2、从而拿到对应的 policy
	// 3、为这些 policy 重新分配

	// 1、找出选择该节点作为网关节点的 EgressGateway
	for _, eg := range egList.Items {
		for _, eipStatus := range eg.Status.NodeList {
			if nodeName == eipStatus.Name {
				err := r.deleteNodeFromEG(ctx, nodeName, eg)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

// 将该节点从 EG 中删除
func (r egnReconciler) deleteNodeFromEG(ctx context.Context, nodeName string, eg egress.EgressGateway) error {
	// 拿到需要重新分配的 policy
	policies, isExist := GetEIPStatusByNode(nodeName, eg)

	if isExist {
		perNodeListMap := make(map[string]egress.EgressIPStatus, 0)
		for _, item := range eg.Status.NodeList {
			if nodeName != item.Name {
				perNodeListMap[item.Name] = item
			}
		}

		// 重新分配网关节点
		for _, policy := range policies {
			err := r.reAllocatorPolicy(ctx, policy, &eg, perNodeListMap)
			if err != nil {
				r.log.Sugar().Errorf("reallocator Failed to reassign a gateway node for EgressGatewayPolicy %v: %v", policy, err)
				return err
			}
		}

		var perNodeList []egress.EgressIPStatus
		for _, node := range perNodeListMap {
			perNodeList = append(perNodeList, node)
		}

		eg.Status.NodeList = perNodeList
		r.log.Sugar().Debugf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
		err := r.client.Status().Update(ctx, &eg)
		if err != nil {
			r.log.Sugar().Errorf("update egress gateway status\n%s", mustMarshalJson(eg.Status))
			return err
		}
	}

	return nil
}

// 重新为 policy 选择网关节点
func (r egnReconciler) reAllocatorPolicy(ctx context.Context, policy egress.Policy, eg *egress.EgressGateway, nodeListMap map[string]egress.EgressIPStatus) error {
	egp := &egress.EgressGatewayPolicy{}
	err := r.client.Get(ctx, types.NamespacedName{Namespace: policy.Namespace, Name: policy.Name}, egp)
	if err != nil {
		return err
	}

	// 选择 网关节点
	perNode, err := r.allocatorNode("rr", nodeListMap)
	if err != nil {
		return err
	}

	// 分配 EIP
	ipv4, ipv6, err := r.allocatorEIP("", perNode, *egp, *eg)
	if err != nil {
		return err
	}

	// 更新 EIP及policy
	err = setEipStatus(ipv4, ipv6, perNode, policy, nodeListMap)
	if err != nil {
		return err
	}

	return nil
}

// 分配网关节点
func (r egnReconciler) allocatorNode(selNodePolicy string, nodeListMap map[string]egress.EgressIPStatus) (string, error) {

	if len(nodeListMap) == 0 {
		err := fmt.Errorf("nodeList is empty")
		return "", err
	}
	selNodePolicy = "rr"

	var perNode string
	perNodePolicyNum := 0
	i := 0
	for _, node := range nodeListMap {
		policyNum := 0
		for _, eip := range node.Eips {
			policyNum += len(eip.Policies)
		}

		if i == 0 {
			i++
			perNode = node.Name
			perNodePolicyNum = policyNum
		} else if policyNum <= perNodePolicyNum {
			perNode = node.Name
			perNodePolicyNum = policyNum
		}
	}

	return perNode, nil
}

// 分配 EIP
func (r egnReconciler) allocatorEIP(selEipLolicy string, nodeName string, egp egress.EgressGatewayPolicy, eg egress.EgressGateway) (string, string, error) {

	selEipLolicy = "rr"

	if egp.Spec.EgressIP.UseNodeIP {
		return "", "", nil
	}

	var perIpv4 string
	var perIpv6 string

	if r.config.FileConfig.EnableIPv4 {
		var useIpv4s []net.IP
		var useIpv4sByNode []net.IP

		ipv4Ranges, _ := utils.MergeIPRanges(constant.IPv4, eg.Spec.Ranges.IPv4)

		perIpv4 = egp.Spec.EgressIP.IPv4
		if len(perIpv4) != 0 {
			result, err := utils.IsIPIncludedRange(constant.IPv4, perIpv4, ipv4Ranges)
			if err != nil {
				return "", "", err
			}
			if !result {
				return "", "", fmt.Errorf("%v is not within the EIP range of EgressGateway %v", perIpv4, eg.Name)
			}
		} else {
			// allocator ip
			// 1、拿到已分配的 IP，计算出出未分配的IP
			// 2、如果所有 IP 已使用，则从该节点的 EIP 中随机分配一个

			for _, node := range eg.Status.NodeList {
				for _, eip := range node.Eips {
					if len(eip.IPv4) != 0 {
						useIpv4s = append(useIpv4s, net.IP(eip.IPv4))
					}
				}
			}

			ipv4s, _ := utils.ParseIPRanges(constant.IPv4, ipv4Ranges)
			freeIpv4s := utils.IPsDiffSet(ipv4s, useIpv4s, false)

			if len(freeIpv4s) == 0 {
				for _, node := range eg.Status.NodeList {
					if node.Name == nodeName {
						for _, eip := range node.Eips {
							if len(eip.IPv4) != 0 {
								useIpv4sByNode = append(useIpv4s, net.IP(eip.IPv4))
							}
						}
					}
				}

				rand.Seed(time.Now().UnixNano())
				perIpv4 = useIpv4sByNode[rand.Intn(len(useIpv4sByNode))].String()
			} else {
				rand.Seed(time.Now().UnixNano())
				perIpv4 = freeIpv4s[rand.Intn(len(freeIpv4s))].String()
			}
		}
	}

	if r.config.FileConfig.EnableIPv6 {
		if len(perIpv4) != 0 {
			return perIpv4, GetEipByIP(perIpv4, eg).IPv6, nil
		}

		var useIpv6s []net.IP
		var useIpv6sByNode []net.IP

		ipv6Ranges, _ := utils.MergeIPRanges(constant.IPv6, eg.Spec.Ranges.IPv6)

		perIpv6 = egp.Spec.EgressIP.IPv6
		if len(perIpv6) != 0 {
			result, err := utils.IsIPIncludedRange(constant.IPv4, perIpv4, ipv6Ranges)
			if err != nil {
				return "", "", err
			}
			if !result {
				return "", "", fmt.Errorf("%v is not within the EIP range of EgressGateway %v", perIpv6, eg.Name)
			}
		} else {
			for _, node := range eg.Status.NodeList {
				for _, eip := range node.Eips {
					if len(eip.IPv6) != 0 {
						useIpv6s = append(useIpv6s, net.IP(eip.IPv6))
					}
				}
			}

			ipv6s, _ := utils.ParseIPRanges(constant.IPv6, ipv6Ranges)
			freeIpv6s := utils.IPsDiffSet(ipv6s, useIpv6s, false)

			if len(freeIpv6s) == 0 {
				for _, node := range eg.Status.NodeList {
					if node.Name == nodeName {
						for _, eip := range node.Eips {
							if len(eip.IPv6) != 0 {
								useIpv6sByNode = append(useIpv6s, net.IP(eip.IPv6))
							}
						}
					}
				}

				rand.Seed(time.Now().UnixNano())
				perIpv6 = useIpv6sByNode[rand.Intn(len(useIpv6sByNode))].String()
			} else {
				rand.Seed(time.Now().UnixNano())
				perIpv6 = freeIpv6s[rand.Intn(len(freeIpv6s))].String()
			}
		}
	}

	return perIpv4, perIpv6, nil
}

func NewEgressGatewayController(mgr manager.Manager, log *zap.Logger, cfg *config.Config) error {
	if log == nil {
		return fmt.Errorf("log can not be nil")
	}
	if cfg == nil {
		return fmt.Errorf("cfg can not be nil")
	}
	r := &egnReconciler{
		client: mgr.GetClient(),
		log:    log,
		config: cfg,
	}

	c, err := controller.New("egressGateway", mgr,
		controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	if err := c.Watch(&source.Kind{Type: &egress.EgressGateway{}},
		handler.EnqueueRequestsFromMapFunc(utils.KindToMapFlat("EgressGateway"))); err != nil {
		return fmt.Errorf("failed to watch EgressGateway: %w", err)
	}

	if err := c.Watch(&source.Kind{Type: &corev1.Node{}},
		handler.EnqueueRequestsFromMapFunc(utils.KindToMapFlat("Node"))); err != nil {
		return fmt.Errorf("failed to watch Node: %w", err)
	}

	if err := c.Watch(&source.Kind{Type: &egress.EgressGatewayPolicy{}},
		handler.EnqueueRequestsFromMapFunc(utils.KindToMapFlat("EgressGatewayPolicy"))); err != nil {
		return fmt.Errorf("failed to watch Node: %w", err)
	}

	if err := c.Watch(&source.Kind{Type: &egress.EgressNode{}},
		handler.EnqueueRequestsFromMapFunc(utils.KindToMapFlat("EgressNode"))); err != nil {
		return fmt.Errorf("failed to watch Node: %w", err)
	}

	return nil
}

// 从 EG 获取单个 EIP 的信息
func GetEipByIP(ipv4 string, eg egress.EgressGateway) egress.Eips {
	var eipInfo egress.Eips
	for _, node := range eg.Status.NodeList {
		for _, eip := range node.Eips {
			if eip.IPv4 == ipv4 {
				eipInfo = eip
			}
		}
	}

	return eipInfo
}

// 更新 EG 的 EIP信息
func setEipStatus(ipv4, ipv6 string, nodeName string, policy egress.Policy, nodeListMap map[string]egress.EgressIPStatus) error {
	eipStatus, ok := nodeListMap[nodeName]
	if !ok {
		return fmt.Errorf("the %v node is not a gateway node", nodeName)
	}
	isExist := false
	newEipStatus := egress.EgressIPStatus{}

	for _, eip := range eipStatus.Eips {
		if ipv4 == eip.IPv4 {
			eip.Policies = append(eip.Policies, policy)

			isExist = true
		}
		newEipStatus.Eips = append(newEipStatus.Eips, eip)
		break
	}

	if !isExist {
		newEip := egress.Eips{}
		newEip.IPv4 = ipv4
		newEip.IPv6 = ipv6
		newEip.Policies = append(newEip.Policies, policy)
		eipStatus.Eips = append(eipStatus.Eips, newEip)
		nodeListMap[nodeName] = eipStatus
	} else {
		nodeListMap[nodeName] = newEipStatus
	}

	return nil
}

func mustMarshalJson(obj interface{}) string {
	raw, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(raw)
}

// 获取 EG 中某个节点的 policy 信息
func GetEIPStatusByNode(nodeName string, eg egress.EgressGateway) ([]egress.Policy, bool) {

	var eipStatus egress.EgressIPStatus
	var policies []egress.Policy
	isExist := false
	for _, node := range eg.Status.NodeList {
		if node.Name == nodeName {
			eipStatus = node
			isExist = true
		}
	}

	if isExist {
		for _, eip := range eipStatus.Eips {
			policies = append(policies, eip.Policies...)
		}
	}

	return policies, isExist
}

// 从 EG 获取某个 policy 的网关节点信息
func GetEIPStatusByPolicy(policy egress.Policy, eg egress.EgressGateway) (egress.EgressIPStatus, bool) {
	var eipStatus egress.EgressIPStatus
	isExist := false

	for _, item := range eg.Status.NodeList {
		for _, eip := range item.Eips {
			for _, p := range eip.Policies {
				if p == policy {
					eipStatus = item
					isExist = true
				}
			}
		}
	}

	return eipStatus, isExist
}

// 从一个 EG 中删除某个 policy
func DeletePlocyFromEG(policy egress.Policy, eg *egress.EgressGateway) {
	var policies []egress.Policy
	var eips []egress.Eips

	for i, node := range eg.Status.NodeList {
		for j, eip := range node.Eips {
			for k, item := range eip.Policies {
				if item == policy {
					policies = append(eip.Policies[:k], eip.Policies[k+1:]...)

					if len(policies) == 0 {
						// Release EIP
						for x, e := range node.Eips {
							if eip.IPv4 == e.IPv4 || eip.IPv6 == e.IPv6 {
								eips = append(node.Eips[:x], node.Eips[x+1:]...)
								break
							}
						}
						eg.Status.NodeList[i].Eips = eips
					} else {
						eg.Status.NodeList[i].Eips[j].Policies = policies
					}
					goto breakHere
				}
			}
		}
	}
breakHere:
	return
}
