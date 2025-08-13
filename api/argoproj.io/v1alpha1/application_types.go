/*
Copyright 2025 coindcx.

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

// DestinationSpec defines the destination configuration
type DestinationSpec struct {
	// Name of the destination
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// SourceSpec defines the source configuration
type SourceSpec struct {
	// RepoURL is the URL to the repository
	// +kubebuilder:validation:Required
	RepoURL string `json:"repoURL"`
}

// SyncSpec defines the sync status
type SyncSpec struct {
	// Status of the sync
	Status string `json:"status"`
}

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// Destination defines where to deploy
	// +optional
	Destination DestinationSpec `json:"destination,omitempty"`

	// Source defines the source configuration
	// +optional
	Source SourceSpec `json:"source,omitempty"`
}

// ApplicationStatus defines the observed state of Application.
type ApplicationStatus struct {
	// Sync status information
	// +optional
	Sync SyncSpec `json:"sync,omitempty"`

	// Conditions represent the latest available observations
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Application
	// +required
	Spec ApplicationSpec `json:"spec"`

	// status defines the observed state of Application
	// +optional
	Status ApplicationStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
