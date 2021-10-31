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

package dispatcher

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"

	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"

	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
	"knative.dev/eventing/pkg/channel"
	"knative.dev/eventing/pkg/channel/fanout"
	"knative.dev/eventing/pkg/channel/multichannelfanout"
	"knative.dev/eventing/pkg/kncloudevents"

	v1alpha1 "github.com/odacremolbap/loudvents/pkg/apis/messaging/v1alpha1"
	messagingv1alpha1 "github.com/odacremolbap/loudvents/pkg/client/generated/clientset/internalclientset/typed/messaging/v1alpha1"
	reconcilerv1alpha1 "github.com/odacremolbap/loudvents/pkg/client/generated/injection/reconciler/messaging/v1alpha1/loudventschannel"
)

// Reconciler reconciles LodVent Channels.
type Reconciler struct {
	multiChannelMessageHandler multichannelfanout.MultiChannelMessageHandler
	reporter                   channel.StatsReporter
	messagingClientSet         messagingv1alpha1.MessagingV1alpha1Interface
}

// Check the interfaces Reconciler should implement
var (
	_ reconcilerv1alpha1.Interface         = (*Reconciler)(nil)
	_ reconcilerv1alpha1.ReadOnlyInterface = (*Reconciler)(nil)
)

// ReconcileKind implements loudventchannel.Interface.
func (r *Reconciler) ReconcileKind(ctx context.Context, lvc *v1alpha1.LoudVentsChannel) reconciler.Event {
	if err := r.reconcile(ctx, lvc); err != nil {
		return err
	}

	// Then patch the subscribers to reflect that they are now ready to go
	return r.patchSubscriberStatus(ctx, lvc)
}

// ObserveKind implements loudventchannel.ReadOnlyInterface.
func (r *Reconciler) ObserveKind(ctx context.Context, lvc *v1alpha1.LoudVentsChannel) reconciler.Event {
	return r.reconcile(ctx, lvc)
}

func (r *Reconciler) reconcile(ctx context.Context, lvc *v1alpha1.LoudVentsChannel) reconciler.Event {
	logging.FromContext(ctx).Infow("Reconciling", zap.Any("LoudVentChannel", lvc))

	if !lvc.IsReady() {
		logging.FromContext(ctx).Debug("lvc is not ready, skipping")
		return nil
	}

	config, err := newConfigForLoudVentChannel(lvc)
	if err != nil {
		logging.FromContext(ctx).Error("Error creating config for loudvent channels", zap.Error(err))
		return err
	}

	// First grab the MultiChannelFanoutMessage handler
	handler := r.multiChannelMessageHandler.GetChannelHandler(config.HostName)
	if handler == nil {
		// No handler yet, create one.
		fanoutHandler, err := fanout.NewFanoutMessageHandler(
			logging.FromContext(ctx).Desugar(),
			channel.NewMessageDispatcher(logging.FromContext(ctx).Desugar()),
			config.FanoutConfig,
			r.reporter,
		)
		if err != nil {
			logging.FromContext(ctx).Error("Failed to create a new fanout.MessageHandler", err)
			return err
		}
		r.multiChannelMessageHandler.SetChannelHandler(config.HostName, fanoutHandler)
	} else {
		// Just update the config if necessary.
		haveSubs := handler.GetSubscriptions(ctx)

		// Ignore the closures, we stash the values that we can tell from if the values have actually changed.
		if diff := cmp.Diff(config.FanoutConfig.Subscriptions, haveSubs, cmpopts.IgnoreFields(kncloudevents.RetryConfig{}, "Backoff", "CheckRetry")); diff != "" {
			logging.FromContext(ctx).Info("Updating fanout config: ", zap.String("Diff", diff))
			handler.SetSubscriptions(ctx, config.FanoutConfig.Subscriptions)
		}
	}

	return nil
}

func (r *Reconciler) patchSubscriberStatus(ctx context.Context, lvc *v1alpha1.LoudVentsChannel) error {
	after := lvc.DeepCopy()

	after.Status.Subscribers = make([]eventingduckv1.SubscriberStatus, 0)
	for _, sub := range lvc.Spec.Subscribers {
		after.Status.Subscribers = append(after.Status.Subscribers, eventingduckv1.SubscriberStatus{
			UID:                sub.UID,
			ObservedGeneration: sub.Generation,
			Ready:              corev1.ConditionTrue,
		})
	}
	jsonPatch, err := duck.CreatePatch(lvc, after)
	if err != nil {
		return fmt.Errorf("creating JSON patch: %w", err)
	}
	// If there is nothing to patch, we are good, just return.
	// Empty patch is [], hence we check for that.
	if len(jsonPatch) == 0 {
		return nil
	}

	patch, err := jsonPatch.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshaling JSON patch: %w", err)
	}
	patched, err := r.messagingClientSet.LoudVentsChannels(lvc.Namespace).Patch(ctx, lvc.Name, types.JSONPatchType, patch, metav1.PatchOptions{}, "status")
	if err != nil {
		return fmt.Errorf("Failed patching: %w", err)
	}
	logging.FromContext(ctx).Debugw("Patched resource", zap.Any("patch", patch), zap.Any("patched", patched))
	return nil
}

// newConfigForLoudVentChannel creates a new Config for a single loudvent channel.
func newConfigForLoudVentChannel(lvc *v1alpha1.LoudVentsChannel) (*multichannelfanout.ChannelConfig, error) {
	subs := make([]fanout.Subscription, len(lvc.Spec.Subscribers))

	for i, sub := range lvc.Spec.Subscribers {
		conf, err := fanout.SubscriberSpecToFanoutConfig(sub)
		if err != nil {
			return nil, err
		}
		subs[i] = *conf
	}

	return &multichannelfanout.ChannelConfig{
		Namespace: lvc.Namespace,
		Name:      lvc.Name,
		HostName:  lvc.Status.Address.URL.Host,
		FanoutConfig: fanout.Config{
			AsyncHandler:  true,
			Subscriptions: subs,
		},
	}, nil
}

func (r *Reconciler) deleteFunc(obj interface{}) {
	if obj == nil {
		return
	}
	acc, err := kmeta.DeletionHandlingAccessor(obj)
	if err != nil {
		return
	}
	lvc, ok := acc.(*v1alpha1.LoudVentsChannel)
	if !ok || lvc == nil {
		return
	}
	if lvc.Status.Address != nil && lvc.Status.Address.URL != nil {
		if hostName := lvc.Status.Address.URL.Host; hostName != "" {
			r.multiChannelMessageHandler.DeleteChannelHandler(hostName)
		}
	}
}
