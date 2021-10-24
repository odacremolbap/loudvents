/*
Copyright 2020 The Knative Authors

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LoudVentsChannel is a resource representing a loudvents channel
type LoudVentsChannel struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of the Channel.
	Spec LoudVentsChannelSpec `json:"spec,omitempty"`

	// Status represents the current state of the Channel. This data may be out of
	// date.
	// +optional
	Status LoudVentsChannelStatus `json:"status,omitempty"`
}

var (
	// Check that InMemoryChannel can return its spec untyped.
	_ apis.HasSpec = (*LoudVentsChannel)(nil)

	_ runtime.Object = (*LoudVentsChannel)(nil)

	// Check that we can create OwnerReferences to an InMemoryChannel.
	_ kmeta.OwnerRefable = (*LoudVentsChannel)(nil)

	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*LoudVentsChannel)(nil)
)

// LoudVentsChannelSpec defines which subscribers have expressed interest in
// receiving events from this InMemoryChannel.
// arguments for a Channel.
type LoudVentsChannelSpec struct {
	// Channel conforms to Duck type Channelable.
	eventingduckv1.ChannelableSpec `json:",inline"`
}

// LoudVentsChannelStatus represents the current state of a Channel.
type LoudVentsChannelStatus struct {
	// Channel conforms to Duck type ChannelableStatus.
	eventingduckv1.ChannelableStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LoudVentsChannelList is a collection of in-memory channels.
type LoudVentsChannelList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoudVentsChannel `json:"items"`
}

// GetStatus retrieves the status of the InMemoryChannel. Implements the KRShaped interface.
func (t *LoudVentsChannel) GetStatus() *duckv1.Status {
	return &t.Status.Status
}
