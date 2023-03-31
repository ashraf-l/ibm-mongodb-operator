/*
# Copyright 2021 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

package v1alpha1

import (
	"context"
	"reflect"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type Image struct {
	Tag string `json:"tag,omitempty"`
}

// MongoDBSpec defines the desired state of MongoDB
type MongoDBSpec struct {
	ImageRegistry  string                      `json:"imageRegistry,omitempty"`
	Replicas       int                         `json:"replicas,omitempty"`
	StorageClass   string                      `json:"storageClass,omitempty"`
	InitImage      Image                       `json:"initImage,omitempty"`
	BootstrapImage Image                       `json:"bootstrapImage,omitempty"`
	MetricsImage   Image                       `json:"metricsImage,omitempty"`
	Resources      corev1.ResourceRequirements `json:"resources,omitempty"`
	PVC            MongoDBPVCSpec              `json:"pvc,omitempty"`
}

// MongoDBPVCSpec defines the desired state of the MongoDB PVCs
type MongoDBPVCSpec struct {
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// MongoDBStatus defines the observed state of MongoDB
type MongoDBStatus struct {
	StorageClass string        `json:"storageClass,omitempty"`
	Service      ServiceStatus `json:"service,omitempty"`
}

// +kubebuilder:object:root=true

// MongoDB is the Schema for the mongodbs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=mongodbs,scope=Namespaced
type MongoDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoDBSpec   `json:"spec,omitempty"`
	Status MongoDBStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MongoDBList contains a list of MongoDB
type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDB `json:"items"`
}

type ManagedResourceStatus struct {
	ObjectName string `json:"objectName,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Status     string `json:"status,omitempty"`
}

type ServiceStatus struct {
	ObjectName       string                  `json:"objectName,omitempty"`
	APIVersion       string                  `json:"apiVersion,omitempty"`
	Namespace        string                  `json:"namespace,omitempty"`
	Kind             string                  `json:"kind,omitempty"`
	Status           string                  `json:"status,omitempty"`
	ManagedResources []ManagedResourceStatus `json:"managedResources,omitempty"`
}

func (a *MongoDB) SetService(ctx context.Context, service ServiceStatus, statusClient client.StatusClient, mu sync.Locker) (err error) {
	reqLogger := logf.FromContext(ctx).WithName("SetService")
	mu.Lock()
	defer mu.Unlock()

	updatedServiceStatus := false
	if !reflect.DeepEqual(service, a.Status.Service) {
		a.Status.Service = service
		updatedServiceStatus = true
	}

	if updatedServiceStatus {
		reqLogger.Info("Status has changed; performing update")
		err = statusClient.Status().Update(ctx, a)
	} else {
		reqLogger.Info("Status is the same; skipping update")
	}
	if err != nil {
		reqLogger.Error(err, "Attempt to update failed")
	}
	return nil
}

func init() {
	SchemeBuilder.Register(&MongoDB{}, &MongoDBList{})
}
