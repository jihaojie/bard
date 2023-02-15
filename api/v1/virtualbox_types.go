/*
Copyright 2023 jihaojie.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type BoxType string
type DiskType string
type BoxStatus string

const (
	BoxTypeDefault BoxType = "Default"

	BoxStatusPending BoxStatus = "Pending"
	BoxStatusRunning BoxStatus = "Running"
	BoxStatusError   BoxStatus = "Error"
	BoxStatusUnknown BoxStatus = "Unknown"
	BoxStatusStopped BoxStatus = "Stopped"

	DiskTypeHost DiskType = "Host"
	DiskTypeCeph DiskType = "Ceph"
	DiskTypeNFS  DiskType = "NFS"
)

// VirtualboxSpec defines the desired state of Virtualbox
type VirtualboxSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	BoxType BoxType `json:"boxType"` //虚拟环境类型
	Size    *int32  `json:"size"`    //开几个虚拟机实例
	Memory  string  `json:"memory"`  //内存大小
	CPU     int     `json:"cpu"`     //CPU大小

	DiskType DiskType `json:"diskType"` //存储类型
	DiskSize string   `json:"diskSize"` //存储大小

	//TODO:考虑计费属性字段
}

// VirtualboxStatus defines the observed state of Virtualbox
type VirtualboxStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase BoxStatus `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Virtualbox is the Schema for the virtualboxes API
type Virtualbox struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualboxSpec   `json:"spec,omitempty"`
	Status VirtualboxStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VirtualboxList contains a list of Virtualbox
type VirtualboxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Virtualbox `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Virtualbox{}, &VirtualboxList{})
}
