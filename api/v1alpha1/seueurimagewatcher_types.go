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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SeueurImageWatcherSpec defines the desired state of SeueurImageWatcher
type SeueurImageWatcherSpec struct {
	// Image is the full image name without tag (e.g. ghcr.io/angga/backend)
	Image string `json:"image"`
	
	// Tag to monitor (e.g. production, latest)
	Tag string `json:"tag"`
	
	// UpdatePolicy: static, semver, or none
	UpdatePolicy string `json:"updatePolicy"`
	
	// AuthRef references a Docker registry secret
	// +optional
	AuthRef *SecretReference `json:"authRef,omitempty"`
	
	// WebhookRef references a webhook configuration
	// +optional
	WebhookRef *LocalObjectReference `json:"webhookRef,omitempty"`
	
	// TargetRef specifies the workload to update
	TargetRef WorkloadReference `json:"targetRef"`
	
	// IntervalSeconds defines how often to check for updates (default: 600)
	// +optional
	IntervalSeconds int `json:"intervalSeconds,omitempty"`
}

// SecretReference references a Kubernetes secret
type SecretReference struct {
	Name string `json:"name"`
}

// LocalObjectReference references another object in the same namespace
type LocalObjectReference struct {
	Name string `json:"name"`
}

// WorkloadReference references a Kubernetes workload
type WorkloadReference struct {
	// Kind: Deployment, StatefulSet, or DaemonSet
	Kind string `json:"kind"`
	
	// Name of the workload
	Name string `json:"name"`
	
	// Namespace of the workload
	Namespace string `json:"namespace"`
}

// SeueurImageWatcherStatus defines the observed state of SeueurImageWatcher
type SeueurImageWatcherStatus struct {
	// LastDigest is the last observed image digest
	// +optional
	LastDigest string `json:"lastDigest,omitempty"`
	
	// LastChecked is when the image was last checked
	// +optional
	LastChecked metav1.Time `json:"lastChecked,omitempty"`
	
	// LastUpdated is when the workload was last updated
	// +optional
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`
	
	// LastNotified is when a notification was last sent
	// +optional
	LastNotified metav1.Time `json:"lastNotified,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Tag",type="string",JSONPath=".spec.tag"
//+kubebuilder:printcolumn:name="Policy",type="string",JSONPath=".spec.updatePolicy"
//+kubebuilder:printcolumn:name="Last Digest",type="string",JSONPath=".status.lastDigest"
//+kubebuilder:printcolumn:name="Last Checked",type="string",JSONPath=".status.lastChecked"

// SeueurImageWatcher is the Schema for the seueurimagewatchers API
type SeueurImageWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SeueurImageWatcherSpec   `json:"spec,omitempty"`
	Status SeueurImageWatcherStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SeueurImageWatcherList contains a list of SeueurImageWatcher
type SeueurImageWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SeueurImageWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SeueurImageWatcher{}, &SeueurImageWatcherList{})
}