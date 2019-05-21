package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RestService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RestServiceSpec   `json:"spec"`
	Status RestServiceStatus `json:"status"`
}

// RestServiceSpec is the spec for a Foo resource
type RestServiceSpec struct {
	DeploymentName string `json:"deploymentName"`
	NameSpace      string `json:"namespace"`
	Image          string `json:"image"`
	Replicas       *int32 `json:"replicas"`
}

// FooStatus is the status for a Foo resource
type RestServiceStatus struct {
	AvailableReplicas int32    `json:"availableReplicas"`
	ActionHistory     []string `json:"actionHistory"`
	VerifyCmd         string   `json:"verifyCommand"`
	ServiceIP         string   `json:"serviceIP"`
	ServicePort       string   `json:"servicePort"`
	Status            string   `json:"status"`
}

type RestServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RestService `json:"items"`
}
