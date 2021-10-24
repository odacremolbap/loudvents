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

package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	v1 "knative.dev/pkg/apis/duck/v1"
)

var imcCondSet = apis.NewLivingConditionSet(
	LoudVentsChannelConditionDispatcherReady,
	LoudVentsChannelConditionServiceReady,
	LoudVentsChannelConditionEndpointsReady,
	LoudVentsChannelConditionAddressable,
	LoudVentsChannelConditionChannelServiceReady,
	LoudVentsChannelConditionDeadLetterSinkResolved,
)

const (
	// LoudVentsChannelConditionReady has status True when all subconditions below have been set to True.
	LoudVentsChannelConditionReady = apis.ConditionReady

	// LoudVentsChannelConditionDispatcherReady has status True when a Dispatcher deployment is ready
	// Keyed off appsv1.DeploymentAvailable, which means minimum available replicas required are up
	// and running for at least minReadySeconds.
	LoudVentsChannelConditionDispatcherReady apis.ConditionType = "DispatcherReady"

	// LoudVentsChannelConditionServiceReady has status True when a k8s Service is ready. This
	// basically just means it exists because there's no meaningful status in Service. See Endpoints
	// below.
	LoudVentsChannelConditionServiceReady apis.ConditionType = "ServiceReady"

	// LoudVentsChannelConditionEndpointsReady has status True when a k8s Service Endpoints are backed
	// by at least one endpoint.
	LoudVentsChannelConditionEndpointsReady apis.ConditionType = "EndpointsReady"

	// LoudVentsChannelConditionAddressable has status true when this LoudVentsChannel meets
	// the Addressable contract and has a non-empty hostname.
	LoudVentsChannelConditionAddressable apis.ConditionType = "Addressable"

	// LoudVentsChannelConditionServiceReady has status True when a k8s Service representing the channel is ready.
	// Because this uses ExternalName, there are no endpoints to check.
	LoudVentsChannelConditionChannelServiceReady apis.ConditionType = "ChannelServiceReady"

	// LoudVentsChannelConditionDeadLetterSinkResolved has status True when there is a Dead Letter Sink ref or URI
	// defined in the Spec.Delivery, is a valid destination and its correctly resolved into a valid URI
	LoudVentsChannelConditionDeadLetterSinkResolved apis.ConditionType = "DeadLetterSinkResolved"
)

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*LoudVentsChannel) GetConditionSet() apis.ConditionSet {
	return imcCondSet
}

// GetGroupVersionKind returns GroupVersionKind for LoudVentsChannels
func (*LoudVentsChannel) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("LoudVentsChannel")
}

// GetUntypedSpec returns the spec of the LoudVentsChannel.
func (i *LoudVentsChannel) GetUntypedSpec() interface{} {
	return i.Spec
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (imcs *LoudVentsChannelStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return imcCondSet.Manage(imcs).GetCondition(t)
}

// IsReady returns true if the Status condition LoudVentsChannelConditionReady
// is true and the latest spec has been observed.
func (imc *LoudVentsChannel) IsReady() bool {
	imcs := imc.Status
	return imcs.ObservedGeneration == imc.Generation &&
		imc.GetConditionSet().Manage(&imcs).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (imcs *LoudVentsChannelStatus) InitializeConditions() {
	imcCondSet.Manage(imcs).InitializeConditions()
}

func (imcs *LoudVentsChannelStatus) SetAddress(url *apis.URL) {
	imcs.Address = &v1.Addressable{URL: url}
	if url != nil {
		imcCondSet.Manage(imcs).MarkTrue(LoudVentsChannelConditionAddressable)
	} else {
		imcCondSet.Manage(imcs).MarkFalse(LoudVentsChannelConditionAddressable, "emptyHostname", "hostname is the empty string")
	}
}

func (imcs *LoudVentsChannelStatus) MarkDispatcherFailed(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkFalse(LoudVentsChannelConditionDispatcherReady, reason, messageFormat, messageA...)
}

func (imcs *LoudVentsChannelStatus) MarkDispatcherUnknown(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkUnknown(LoudVentsChannelConditionDispatcherReady, reason, messageFormat, messageA...)
}

// TODO: Unify this with the ones from Eventing. Say: Broker, Trigger.
func (imcs *LoudVentsChannelStatus) PropagateDispatcherStatus(ds *appsv1.DeploymentStatus) {
	for _, cond := range ds.Conditions {
		if cond.Type == appsv1.DeploymentAvailable {
			if cond.Status == corev1.ConditionTrue {
				imcCondSet.Manage(imcs).MarkTrue(LoudVentsChannelConditionDispatcherReady)
			} else if cond.Status == corev1.ConditionFalse {
				imcs.MarkDispatcherFailed("DispatcherDeploymentFalse", "The status of Dispatcher Deployment is False: %s : %s", cond.Reason, cond.Message)
			} else if cond.Status == corev1.ConditionUnknown {
				imcs.MarkDispatcherUnknown("DispatcherDeploymentUnknown", "The status of Dispatcher Deployment is Unknown: %s : %s", cond.Reason, cond.Message)
			}
		}
	}
}

func (imcs *LoudVentsChannelStatus) MarkServiceFailed(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkFalse(LoudVentsChannelConditionServiceReady, reason, messageFormat, messageA...)
}

func (imcs *LoudVentsChannelStatus) MarkServiceUnknown(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkUnknown(LoudVentsChannelConditionServiceReady, reason, messageFormat, messageA...)
}

func (imcs *LoudVentsChannelStatus) MarkServiceTrue() {
	imcCondSet.Manage(imcs).MarkTrue(LoudVentsChannelConditionServiceReady)
}

func (imcs *LoudVentsChannelStatus) MarkChannelServiceFailed(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkFalse(LoudVentsChannelConditionChannelServiceReady, reason, messageFormat, messageA...)
}

func (imcs *LoudVentsChannelStatus) MarkChannelServiceUnknown(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkUnknown(LoudVentsChannelConditionChannelServiceReady, reason, messageFormat, messageA...)
}

func (imcs *LoudVentsChannelStatus) MarkChannelServiceTrue() {
	imcCondSet.Manage(imcs).MarkTrue(LoudVentsChannelConditionChannelServiceReady)
}

func (imcs *LoudVentsChannelStatus) MarkEndpointsFailed(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkFalse(LoudVentsChannelConditionEndpointsReady, reason, messageFormat, messageA...)
}

func (imcs *LoudVentsChannelStatus) MarkEndpointsUnknown(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkUnknown(LoudVentsChannelConditionEndpointsReady, reason, messageFormat, messageA...)
}

func (imcs *LoudVentsChannelStatus) MarkEndpointsTrue() {
	imcCondSet.Manage(imcs).MarkTrue(LoudVentsChannelConditionEndpointsReady)
}

func (imcs *LoudVentsChannelStatus) MarkDeadLetterSinkResolvedSucceeded(deadLetterSinkURI *apis.URL) {
	imcs.DeliveryStatus.DeadLetterSinkURI = deadLetterSinkURI
	imcCondSet.Manage(imcs).MarkTrue(LoudVentsChannelConditionDeadLetterSinkResolved)
}

func (imcs *LoudVentsChannelStatus) MarkDeadLetterSinkNotConfigured() {
	imcs.DeadLetterSinkURI = nil
	imcCondSet.Manage(imcs).MarkTrueWithReason(LoudVentsChannelConditionDeadLetterSinkResolved, "DeadLetterSinkNotConfigured", "No dead letter sink is configured.")
}

func (imcs *LoudVentsChannelStatus) MarkDeadLetterSinkResolvedFailed(reason, messageFormat string, messageA ...interface{}) {
	imcs.DeadLetterSinkURI = nil
	imcCondSet.Manage(imcs).MarkFalse(LoudVentsChannelConditionDeadLetterSinkResolved, reason, messageFormat, messageA...)
}
