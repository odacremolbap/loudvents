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
	"fmt"
	"net/http"
	"sync"

	messagingv1alpha1 "github.com/odacremolbap/loudvents/pkg/apis/messaging/v1alpha1"
	messaginglistersv1alpha1 "github.com/odacremolbap/loudvents/pkg/client/generated/listers/messaging/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"knative.dev/eventing/pkg/channel/multichannelfanout"
)

const (
	readinessProbeReady    = http.StatusNoContent
	readinessProbeNotReady = http.StatusServiceUnavailable
	readinessProbeError    = http.StatusInternalServerError
)

// ReadinessChecker can assert the readiness of a component.
type ReadinessChecker interface {
	IsReady() (bool, error)
}

// DispatcherReadyChecker asserts the readiness of a dispatcher for loudvent Channels.
type DispatcherReadyChecker struct {
	// Allows listing observed Channels.
	chLister messaginglistersv1alpha1.LoudVentsChannelLister

	// Allows listing/counting the handlers which have already been registered.
	chMsgHandler multichannelfanout.MultiChannelMessageHandler

	// Allows safe concurrent read/write of 'isReady'.
	sync.Mutex

	// Minor perf tweak, bypass check if we already observed readiness once.
	isReady bool
}

// IsReady implements ReadinessChecker.
// It checks whether the dispatcher has registered a handler for all observed loudvent Channels.
func (c *DispatcherReadyChecker) IsReady() (bool, error) {
	c.Lock()
	defer c.Unlock()

	// readiness already observed at an earlier point in time, short-circuit the check
	if c.isReady {
		return true, nil
	}

	channels, err := c.chLister.List(labels.Everything())
	if err != nil {
		return false, fmt.Errorf("listing cached LoudVentsChannels: %w", err)
	}

	readyChannels := make([]*messagingv1alpha1.LoudVentsChannel, 0, len(channels))
	for _, channel := range channels {
		if channel.IsReady() {
			readyChannels = append(readyChannels, channel)
		}
	}

	c.isReady = c.chMsgHandler.CountChannelHandlers() >= len(readyChannels)
	return c.isReady, nil
}

// readinessCheckerHTTPHandler returns a http.Handler which executes the given ReadinessChecker.
func readinessCheckerHTTPHandler(c ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isReady, err := c.IsReady()
		switch {
		case err != nil:
			w.WriteHeader(readinessProbeError)
		case isReady:
			w.WriteHeader(readinessProbeReady)
		default:
			w.WriteHeader(readinessProbeNotReady)
		}
	}
}
