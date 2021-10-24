//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoudVentsChannel) DeepCopyInto(out *LoudVentsChannel) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoudVentsChannel.
func (in *LoudVentsChannel) DeepCopy() *LoudVentsChannel {
	if in == nil {
		return nil
	}
	out := new(LoudVentsChannel)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LoudVentsChannel) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoudVentsChannelList) DeepCopyInto(out *LoudVentsChannelList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LoudVentsChannel, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoudVentsChannelList.
func (in *LoudVentsChannelList) DeepCopy() *LoudVentsChannelList {
	if in == nil {
		return nil
	}
	out := new(LoudVentsChannelList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LoudVentsChannelList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoudVentsChannelSpec) DeepCopyInto(out *LoudVentsChannelSpec) {
	*out = *in
	in.ChannelableSpec.DeepCopyInto(&out.ChannelableSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoudVentsChannelSpec.
func (in *LoudVentsChannelSpec) DeepCopy() *LoudVentsChannelSpec {
	if in == nil {
		return nil
	}
	out := new(LoudVentsChannelSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoudVentsChannelStatus) DeepCopyInto(out *LoudVentsChannelStatus) {
	*out = *in
	in.ChannelableStatus.DeepCopyInto(&out.ChannelableStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoudVentsChannelStatus.
func (in *LoudVentsChannelStatus) DeepCopy() *LoudVentsChannelStatus {
	if in == nil {
		return nil
	}
	out := new(LoudVentsChannelStatus)
	in.DeepCopyInto(out)
	return out
}
