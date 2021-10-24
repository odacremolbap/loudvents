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

package controller

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"go.uber.org/zap"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	rbacv1listers "k8s.io/client-go/listers/rbac/v1"
	pkgreconciler "knative.dev/pkg/reconciler"

	v1alpha1 "github.com/odacremolbap/loudvents/pkg/apis/messaging/v1alpha1"
	loudventschannelreconciler "github.com/odacremolbap/loudvents/pkg/client/generated/injection/reconciler/messaging/v1alpha1/loudventschannel"
	"github.com/odacremolbap/loudvents/pkg/reconciler/loudvents/controller/config"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/resolver"
)

type Reconciler struct {
	kubeClientSet kubernetes.Interface

	systemNamespace      string
	dispatcherImage      string
	deploymentLister     appsv1listers.DeploymentLister
	serviceLister        corev1listers.ServiceLister
	endpointsLister      corev1listers.EndpointsLister
	serviceAccountLister corev1listers.ServiceAccountLister
	roleBindingLister    rbacv1listers.RoleBindingLister

	eventDispatcherConfigStore *config.EventDispatcherConfigStore

	uriResolver *resolver.URIResolver
}

// Check that our Reconciler implements Interface
var _ loudventschannelreconciler.Interface = (*Reconciler)(nil)

func (r *Reconciler) ReconcileKind(ctx context.Context, lvc *v1alpha1.LoudVentsChannel) pkgreconciler.Event {
	logging.FromContext(ctx).Infow("Reconciling", zap.Any("LoadVentsChannel", lvc))
	return nil
}
