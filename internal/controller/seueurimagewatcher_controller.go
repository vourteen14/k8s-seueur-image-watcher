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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	vrtn14v1alpha1 "github.com/vourteen14/k8s-seueur-image-watcher/api/v1alpha1"
)

const (
	finalizerName = "seueurimagewatcher.vrtn14.sr/finalizer"
)

// SeueurImageWatcherReconciler reconciles a SeueurImageWatcher object
type SeueurImageWatcherReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=vrtn14.sr,resources=seueurimagewatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vrtn14.sr,resources=seueurimagewatchers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vrtn14.sr,resources=seueurimagewatchers/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets;daemonsets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get

func (r *SeueurImageWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("seueurimagewatcher", req.NamespacedName)

	// Fetch the SeueurImageWatcher instance
	watcher := &vrtn14v1alpha1.SeueurImageWatcher{}
	if err := r.Get(ctx, req.NamespacedName, watcher); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Handle deletion
	if !watcher.ObjectMeta.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(watcher, finalizerName) {
			// Perform cleanup
			controllerutil.RemoveFinalizer(watcher, finalizerName)
			if err := r.Update(ctx, watcher); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(watcher, finalizerName) {
		controllerutil.AddFinalizer(watcher, finalizerName)
		if err := r.Update(ctx, watcher); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Check if it's time to reconcile based on interval
	interval := time.Duration(watcher.Spec.IntervalSeconds) * time.Second
	if watcher.Spec.IntervalSeconds == 0 {
		interval = 600 * time.Second // default
	}

	if time.Since(watcher.Status.LastChecked.Time) < interval {
		nextCheck := watcher.Status.LastChecked.Add(interval)
		return ctrl.Result{RequeueAfter: time.Until(nextCheck)}, nil
	}

	// Check for image updates
	updated, newDigest, err := r.checkImageUpdate(ctx, watcher)
	if err != nil {
		log.Error(err, "Failed to check image update")
		return ctrl.Result{}, err
	}

	// Update status with last checked time
	watcher.Status.LastChecked = metav1.Now()
	if err := r.Status().Update(ctx, watcher); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	if !updated {
		return ctrl.Result{RequeueAfter: interval}, nil
	}

	// Image has been updated - handle the update
	if err := r.handleImageUpdate(ctx, watcher, newDigest); err != nil {
		log.Error(err, "Failed to handle image update")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: interval}, nil
}

func (r *SeueurImageWatcherReconciler) checkImageUpdate(ctx context.Context, watcher *vrtn14v1alpha1.SeueurImageWatcher) (bool, string, error) {
	// TODO: Implement actual image digest check
	// This would involve:
	// 1. Getting registry credentials if authRef is specified
	// 2. Querying the registry for the current digest of the image:tag
	// 3. Comparing with status.lastDigest
	
	// Mock implementation - always return true for testing
	return true, "sha256:mockdigest", nil
}

func (r *SeueurImageWatcherReconciler) handleImageUpdate(ctx context.Context, watcher *vrtn14v1alpha1.SeueurImageWatcher, newDigest string) error {
	// Update status with new digest and update time
	watcher.Status.LastDigest = newDigest
	watcher.Status.LastUpdated = metav1.Now()
	if err := r.Status().Update(ctx, watcher); err != nil {
		return err
	}

	// Send webhook notification if configured
	if watcher.Spec.WebhookRef != nil {
		if err := r.sendWebhookNotification(ctx, watcher); err != nil {
			r.Log.Error(err, "Failed to send webhook notification")
			// Continue even if webhook fails
		} else {
			watcher.Status.LastNotified = metav1.Now()
			if err := r.Status().Update(ctx, watcher); err != nil {
				return err
			}
		}
	}

	// Update the target workload
	switch watcher.Spec.UpdatePolicy {
	case "static":
		return r.triggerRollout(ctx, watcher)
	case "semver":
		return r.updateImageTag(ctx, watcher)
	case "none":
		return nil
	default:
		return fmt.Errorf("unknown update policy: %s", watcher.Spec.UpdatePolicy)
	}
}

func (r *SeueurImageWatcherReconciler) sendWebhookNotification(ctx context.Context, watcher *vrtn14v1alpha1.SeueurImageWatcher) error {
	// TODO: Implement webhook notification
	// 1. Get the SeueurWebhookConfig referenced in webhookRef
	// 2. Parse the template with the current data
	// 3. Send HTTP request to the webhook URL
	
	return nil
}

func (r *SeueurImageWatcherReconciler) triggerRollout(ctx context.Context, watcher *vrtn14v1alpha1.SeueurImageWatcher) error {
	// TODO: Implement rollout trigger for static tags
	// This would typically add an annotation to force a rollout
	// without changing the image tag
	
	return nil
}

func (r *SeueurImageWatcherReconciler) updateImageTag(ctx context.Context, watcher *vrtn14v1alpha1.SeueurImageWatcher) error {
	// TODO: Implement image tag update for semver policy
	// This would update the actual image tag in the workload
	
	return nil
}

func (r *SeueurImageWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vrtn14v1alpha1.SeueurImageWatcher{}).
		Complete(r)
}