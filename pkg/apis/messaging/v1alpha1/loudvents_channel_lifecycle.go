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

var loudventCondSet = apis.NewLivingConditionSet(
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
	return loudventCondSet
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
func (lvcs *LoudVentsChannelStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return loudventCondSet.Manage(lvcs).GetCondition(t)
}

// IsReady returns true if the Status condition LoudVentsChannelConditionReady
// is true and the latest spec has been observed.
func (lvc *LoudVentsChannel) IsReady() bool {
	lvcs := lvc.Status
	return lvcs.ObservedGeneration == lvc.Generation &&
		lvc.GetConditionSet().Manage(&lvcs).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (lvcs *LoudVentsChannelStatus) InitializeConditions() {
	loudventCondSet.Manage(lvcs).InitializeConditions()
}

func (lvcs *LoudVentsChannelStatus) SetAddress(url *apis.URL) {
	lvcs.Address = &v1.Addressable{URL: url}
	if url != nil {
		loudventCondSet.Manage(lvcs).MarkTrue(LoudVentsChannelConditionAddressable)
	} else {
		loudventCondSet.Manage(lvcs).MarkFalse(LoudVentsChannelConditionAddressable, "emptyHostname", "hostname is the empty string")
	}
}

func (lvcs *LoudVentsChannelStatus) MarkDispatcherFailed(reason, messageFormat string, messageA ...interface{}) {
	loudventCondSet.Manage(lvcs).MarkFalse(LoudVentsChannelConditionDispatcherReady, reason, messageFormat, messageA...)
}

func (lvcs *LoudVentsChannelStatus) MarkDispatcherUnknown(reason, messageFormat string, messageA ...interface{}) {
	loudventCondSet.Manage(lvcs).MarkUnknown(LoudVentsChannelConditionDispatcherReady, reason, messageFormat, messageA...)
}

// TODO: Unify this with the ones from Eventing. Say: Broker, Trigger.
func (lvcs *LoudVentsChannelStatus) PropagateDispatcherStatus(ds *appsv1.DeploymentStatus) {
	for _, cond := range ds.Conditions {
		if cond.Type == appsv1.DeploymentAvailable {
			if cond.Status == corev1.ConditionTrue {
				loudventCondSet.Manage(lvcs).MarkTrue(LoudVentsChannelConditionDispatcherReady)
			} else if cond.Status == corev1.ConditionFalse {
				lvcs.MarkDispatcherFailed("DispatcherDeploymentFalse", "The status of Dispatcher Deployment is False: %s : %s", cond.Reason, cond.Message)
			} else if cond.Status == corev1.ConditionUnknown {
				lvcs.MarkDispatcherUnknown("DispatcherDeploymentUnknown", "The status of Dispatcher Deployment is Unknown: %s : %s", cond.Reason, cond.Message)
			}
		}
	}
}

func (lvcs *LoudVentsChannelStatus) MarkServiceFailed(reason, messageFormat string, messageA ...interface{}) {
	loudventCondSet.Manage(lvcs).MarkFalse(LoudVentsChannelConditionServiceReady, reason, messageFormat, messageA...)
}

func (lvcs *LoudVentsChannelStatus) MarkServiceUnknown(reason, messageFormat string, messageA ...interface{}) {
	loudventCondSet.Manage(lvcs).MarkUnknown(LoudVentsChannelConditionServiceReady, reason, messageFormat, messageA...)
}

func (lvcs *LoudVentsChannelStatus) MarkServiceTrue() {
	loudventCondSet.Manage(lvcs).MarkTrue(LoudVentsChannelConditionServiceReady)
}

func (lvcs *LoudVentsChannelStatus) MarkChannelServiceFailed(reason, messageFormat string, messageA ...interface{}) {
	loudventCondSet.Manage(lvcs).MarkFalse(LoudVentsChannelConditionChannelServiceReady, reason, messageFormat, messageA...)
}

func (lvcs *LoudVentsChannelStatus) MarkChannelServiceUnknown(reason, messageFormat string, messageA ...interface{}) {
	loudventCondSet.Manage(lvcs).MarkUnknown(LoudVentsChannelConditionChannelServiceReady, reason, messageFormat, messageA...)
}

func (lvcs *LoudVentsChannelStatus) MarkChannelServiceTrue() {
	loudventCondSet.Manage(lvcs).MarkTrue(LoudVentsChannelConditionChannelServiceReady)
}

func (lvcs *LoudVentsChannelStatus) MarkEndpointsFailed(reason, messageFormat string, messageA ...interface{}) {
	loudventCondSet.Manage(lvcs).MarkFalse(LoudVentsChannelConditionEndpointsReady, reason, messageFormat, messageA...)
}

func (lvcs *LoudVentsChannelStatus) MarkEndpointsUnknown(reason, messageFormat string, messageA ...interface{}) {
	loudventCondSet.Manage(lvcs).MarkUnknown(LoudVentsChannelConditionEndpointsReady, reason, messageFormat, messageA...)
}

func (lvcs *LoudVentsChannelStatus) MarkEndpointsTrue() {
	loudventCondSet.Manage(lvcs).MarkTrue(LoudVentsChannelConditionEndpointsReady)
}

func (lvcs *LoudVentsChannelStatus) MarkDeadLetterSinkResolvedSucceeded(deadLetterSinkURI *apis.URL) {
	lvcs.DeliveryStatus.DeadLetterSinkURI = deadLetterSinkURI
	loudventCondSet.Manage(lvcs).MarkTrue(LoudVentsChannelConditionDeadLetterSinkResolved)
}

func (lvcs *LoudVentsChannelStatus) MarkDeadLetterSinkNotConfigured() {
	lvcs.DeadLetterSinkURI = nil
	loudventCondSet.Manage(lvcs).MarkTrueWithReason(LoudVentsChannelConditionDeadLetterSinkResolved, "DeadLetterSinkNotConfigured", "No dead letter sink is configured.")
}

func (lvcs *LoudVentsChannelStatus) MarkDeadLetterSinkResolvedFailed(reason, messageFormat string, messageA ...interface{}) {
	lvcs.DeadLetterSinkURI = nil
	loudventCondSet.Manage(lvcs).MarkFalse(LoudVentsChannelConditionDeadLetterSinkResolved, reason, messageFormat, messageA...)
}
