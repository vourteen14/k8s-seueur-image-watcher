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

// SeueurWebhookConfigSpec defines the desired state of SeueurWebhookConfig
type SeueurWebhookConfigSpec struct {
	// URL of the webhook endpoint
	URL string `json:"url"`
	
	// HTTP method to use (default: POST)
	// +optional
	Method string `json:"method,omitempty"`
	
	// Headers to include in the request
	// +optional
	Headers map[string]string `json:"headers,omitempty"`
	
	// Template for the request body (Go template format)
	Template string `json:"template"`
}

//+kubebuilder:object:root=true
//+kubebuilder:printcolumn:name="URL",type="string",JSONPath=".spec.url"
//+kubebuilder:printcolumn:name="Method",type="string",JSONPath=".spec.method"

// SeueurWebhookConfig is the Schema for the seueurwebhookconfigs API
type SeueurWebhookConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SeueurWebhookConfigSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// SeueurWebhookConfigList contains a list of SeueurWebhookConfig
type SeueurWebhookConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SeueurWebhookConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SeueurWebhookConfig{}, &SeueurWebhookConfigList{})
}