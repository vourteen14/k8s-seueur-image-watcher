/*
Copyright 2025.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	vrtn14v1alpha1 "github.com/vourteen14/k8s-seueur-image-watcher/api/v1alpha1"
)

// SeueurWebhookConfigReconciler reconciles a SeueurWebhookConfig object
type SeueurWebhookConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=vrtn14.sr,resources=seueurwebhookconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vrtn14.sr,resources=seueurwebhookconfigs/status,verbs=get;update;patch

func (r *SeueurWebhookConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("seueurwebhookconfig", req.NamespacedName)

	// Fetch the SeueurWebhookConfig instance
	config := &vrtn14v1alpha1.SeueurWebhookConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Currently no status to update or actions to take
	return ctrl.Result{}, nil
}

func (r *SeueurWebhookConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vrtn14v1alpha1.SeueurWebhookConfig{}).
		Complete(r)
}