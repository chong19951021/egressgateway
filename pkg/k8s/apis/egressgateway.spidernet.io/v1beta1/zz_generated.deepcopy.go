//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedTo) DeepCopyInto(out *AppliedTo) {
	*out = *in
	if in.PodSelector != nil {
		in, out := &in.PodSelector, &out.PodSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.PodSubnet != nil {
		in, out := &in.PodSubnet, &out.PodSubnet
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedTo.
func (in *AppliedTo) DeepCopy() *AppliedTo {
	if in == nil {
		return nil
	}
	out := new(AppliedTo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressClusterInfo) DeepCopyInto(out *EgressClusterInfo) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressClusterInfo.
func (in *EgressClusterInfo) DeepCopy() *EgressClusterInfo {
	if in == nil {
		return nil
	}
	out := new(EgressClusterInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressClusterInfo) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressClusterInfoList) DeepCopyInto(out *EgressClusterInfoList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EgressClusterInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressClusterInfoList.
func (in *EgressClusterInfoList) DeepCopy() *EgressClusterInfoList {
	if in == nil {
		return nil
	}
	out := new(EgressClusterInfoList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressClusterInfoList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressClusterStatus) DeepCopyInto(out *EgressClusterStatus) {
	*out = *in
	in.EgressIgnoreCIDR.DeepCopyInto(&out.EgressIgnoreCIDR)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressClusterStatus.
func (in *EgressClusterStatus) DeepCopy() *EgressClusterStatus {
	if in == nil {
		return nil
	}
	out := new(EgressClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressClusterStatusSpec) DeepCopyInto(out *EgressClusterStatusSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressClusterStatusSpec.
func (in *EgressClusterStatusSpec) DeepCopy() *EgressClusterStatusSpec {
	if in == nil {
		return nil
	}
	out := new(EgressClusterStatusSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressEndpoint) DeepCopyInto(out *EgressEndpoint) {
	*out = *in
	if in.IPv4List != nil {
		in, out := &in.IPv4List, &out.IPv4List
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.IPv6List != nil {
		in, out := &in.IPv6List, &out.IPv6List
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressEndpoint.
func (in *EgressEndpoint) DeepCopy() *EgressEndpoint {
	if in == nil {
		return nil
	}
	out := new(EgressEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressEndpointSlice) DeepCopyInto(out *EgressEndpointSlice) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressEndpointSlice.
func (in *EgressEndpointSlice) DeepCopy() *EgressEndpointSlice {
	if in == nil {
		return nil
	}
	out := new(EgressEndpointSlice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressEndpointSlice) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressEndpointSliceList) DeepCopyInto(out *EgressEndpointSliceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EgressEndpointSlice, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressEndpointSliceList.
func (in *EgressEndpointSliceList) DeepCopy() *EgressEndpointSliceList {
	if in == nil {
		return nil
	}
	out := new(EgressEndpointSliceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressEndpointSliceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressEndpointSliceSpec) DeepCopyInto(out *EgressEndpointSliceSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressEndpointSliceSpec.
func (in *EgressEndpointSliceSpec) DeepCopy() *EgressEndpointSliceSpec {
	if in == nil {
		return nil
	}
	out := new(EgressEndpointSliceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressEndpointSliceSpecStatus) DeepCopyInto(out *EgressEndpointSliceSpecStatus) {
	*out = *in
	if in.Endpoints != nil {
		in, out := &in.Endpoints, &out.Endpoints
		*out = make([]EgressEndpoint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressEndpointSliceSpecStatus.
func (in *EgressEndpointSliceSpecStatus) DeepCopy() *EgressEndpointSliceSpecStatus {
	if in == nil {
		return nil
	}
	out := new(EgressEndpointSliceSpecStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressGateway) DeepCopyInto(out *EgressGateway) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressGateway.
func (in *EgressGateway) DeepCopy() *EgressGateway {
	if in == nil {
		return nil
	}
	out := new(EgressGateway)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressGateway) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressGatewayList) DeepCopyInto(out *EgressGatewayList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EgressGateway, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressGatewayList.
func (in *EgressGatewayList) DeepCopy() *EgressGatewayList {
	if in == nil {
		return nil
	}
	out := new(EgressGatewayList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressGatewayList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressGatewayPolicy) DeepCopyInto(out *EgressGatewayPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressGatewayPolicy.
func (in *EgressGatewayPolicy) DeepCopy() *EgressGatewayPolicy {
	if in == nil {
		return nil
	}
	out := new(EgressGatewayPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressGatewayPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressGatewayPolicyList) DeepCopyInto(out *EgressGatewayPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EgressGatewayPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressGatewayPolicyList.
func (in *EgressGatewayPolicyList) DeepCopy() *EgressGatewayPolicyList {
	if in == nil {
		return nil
	}
	out := new(EgressGatewayPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressGatewayPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressGatewayPolicySpec) DeepCopyInto(out *EgressGatewayPolicySpec) {
	*out = *in
	out.EgressIP = in.EgressIP
	in.AppliedTo.DeepCopyInto(&out.AppliedTo)
	if in.DestSubnet != nil {
		in, out := &in.DestSubnet, &out.DestSubnet
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressGatewayPolicySpec.
func (in *EgressGatewayPolicySpec) DeepCopy() *EgressGatewayPolicySpec {
	if in == nil {
		return nil
	}
	out := new(EgressGatewayPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressGatewaySpec) DeepCopyInto(out *EgressGatewaySpec) {
	*out = *in
	in.Ranges.DeepCopyInto(&out.Ranges)
	in.NodeSelector.DeepCopyInto(&out.NodeSelector)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressGatewaySpec.
func (in *EgressGatewaySpec) DeepCopy() *EgressGatewaySpec {
	if in == nil {
		return nil
	}
	out := new(EgressGatewaySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressGatewayStatus) DeepCopyInto(out *EgressGatewayStatus) {
	*out = *in
	if in.NodeList != nil {
		in, out := &in.NodeList, &out.NodeList
		*out = make([]EgressIPStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressGatewayStatus.
func (in *EgressGatewayStatus) DeepCopy() *EgressGatewayStatus {
	if in == nil {
		return nil
	}
	out := new(EgressGatewayStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressIP) DeepCopyInto(out *EgressIP) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressIP.
func (in *EgressIP) DeepCopy() *EgressIP {
	if in == nil {
		return nil
	}
	out := new(EgressIP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressIPStatus) DeepCopyInto(out *EgressIPStatus) {
	*out = *in
	if in.Eips != nil {
		in, out := &in.Eips, &out.Eips
		*out = make([]Eips, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressIPStatus.
func (in *EgressIPStatus) DeepCopy() *EgressIPStatus {
	if in == nil {
		return nil
	}
	out := new(EgressIPStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressIgnoreCIDR) DeepCopyInto(out *EgressIgnoreCIDR) {
	*out = *in
	in.NodeIP.DeepCopyInto(&out.NodeIP)
	in.ClusterIP.DeepCopyInto(&out.ClusterIP)
	in.PodCIDR.DeepCopyInto(&out.PodCIDR)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressIgnoreCIDR.
func (in *EgressIgnoreCIDR) DeepCopy() *EgressIgnoreCIDR {
	if in == nil {
		return nil
	}
	out := new(EgressIgnoreCIDR)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressNode) DeepCopyInto(out *EgressNode) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressNode.
func (in *EgressNode) DeepCopy() *EgressNode {
	if in == nil {
		return nil
	}
	out := new(EgressNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressNode) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressNodeList) DeepCopyInto(out *EgressNodeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EgressNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressNodeList.
func (in *EgressNodeList) DeepCopy() *EgressNodeList {
	if in == nil {
		return nil
	}
	out := new(EgressNodeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EgressNodeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressNodeSpec) DeepCopyInto(out *EgressNodeSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressNodeSpec.
func (in *EgressNodeSpec) DeepCopy() *EgressNodeSpec {
	if in == nil {
		return nil
	}
	out := new(EgressNodeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressNodeStatus) DeepCopyInto(out *EgressNodeStatus) {
	*out = *in
	out.Tunnel = in.Tunnel
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressNodeStatus.
func (in *EgressNodeStatus) DeepCopy() *EgressNodeStatus {
	if in == nil {
		return nil
	}
	out := new(EgressNodeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Eips) DeepCopyInto(out *Eips) {
	*out = *in
	if in.Policies != nil {
		in, out := &in.Policies, &out.Policies
		*out = make([]Policy, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Eips.
func (in *Eips) DeepCopy() *Eips {
	if in == nil {
		return nil
	}
	out := new(Eips)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPListPair) DeepCopyInto(out *IPListPair) {
	*out = *in
	if in.IPv4 != nil {
		in, out := &in.IPv4, &out.IPv4
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPListPair.
func (in *IPListPair) DeepCopy() *IPListPair {
	if in == nil {
		return nil
	}
	out := new(IPListPair)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeSelector) DeepCopyInto(out *NodeSelector) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeSelector.
func (in *NodeSelector) DeepCopy() *NodeSelector {
	if in == nil {
		return nil
	}
	out := new(NodeSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parent) DeepCopyInto(out *Parent) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parent.
func (in *Parent) DeepCopy() *Parent {
	if in == nil {
		return nil
	}
	out := new(Parent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Policy) DeepCopyInto(out *Policy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Policy.
func (in *Policy) DeepCopy() *Policy {
	if in == nil {
		return nil
	}
	out := new(Policy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ranges) DeepCopyInto(out *Ranges) {
	*out = *in
	if in.IPv4 != nil {
		in, out := &in.IPv4, &out.IPv4
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ranges.
func (in *Ranges) DeepCopy() *Ranges {
	if in == nil {
		return nil
	}
	out := new(Ranges)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tunnel) DeepCopyInto(out *Tunnel) {
	*out = *in
	out.Parent = in.Parent
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tunnel.
func (in *Tunnel) DeepCopy() *Tunnel {
	if in == nil {
		return nil
	}
	out := new(Tunnel)
	in.DeepCopyInto(out)
	return out
}
