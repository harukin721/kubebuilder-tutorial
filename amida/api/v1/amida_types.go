/*
Copyright 2023.

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

// 一般的にカスタムリソースのSpecはユーザーが記述するもので、システムのあるべき状態をユーザーからコントローラーに伝える用途として利用されます。
// AmidaSpec defines the desired state of Amida
type AmidaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Amidaリソースの定義に結果を格納するためのフィールドを追加します。
	// SELECTS is the list of items to select from
	Selects []string `json:"selects,omitempty"`
	// SELECTCOUNT is the number of items to select
	SelectCount int32 `json:"selectCount,omitempty"`
	// RESULTS is the list to store the selected items
	Results []string `json:"results,omitempty"`
}

// カスタムリソースのStatusはコントローラーが記述するもので、実際のシステムの状態を表現するために利用されます。
// AmidaStatus defines the observed state of Amida
type AmidaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Amida is the Schema for the amida API
type Amida struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AmidaSpec   `json:"spec,omitempty"`
	Status AmidaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AmidaList contains a list of Amida
type AmidaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Amida `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Amida{}, &AmidaList{})
}
