package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//go:generate controller-gen object paths=$GOFILE

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TopologySpec struct {
	metav1.TypeMeta `json:",inline"`
	Links           []Link `json:"links"`
}

type Link struct {
	LocalIntf string `json:"local_intf"`
	LocalIP   string `json:"local_ip"`
	PeerIntf  string `json:"peer_intf"`
	PeerIP    string `json:"peer_ip"`
	PeerPod   string `json:"peer_pod"`
	UID       int    `json:"uid"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Topology struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TopologySpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TopologyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Topology `json:"items"`
}
