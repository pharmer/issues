package digitalocean

import (
	"github.com/digitalocean/godo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DigitalOceanMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:",inline"`

	Region        string   `json:"region,omitempty"`
	Size          string   `json:"size,omitempty"`
	Image         string   `json:"image,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	SSHPublicKeys []string `json:"sshPublicKeys,omitempty"`

	PrivateNetworking bool `json:"private_networking,omitempty"`
	Backups           bool `json:"backups,omitempty"`
	IPv6              bool `json:"ipv6,omitempty"`
	Monitoring        bool `json:"monitoring,omitempty"`
}

type DigitalOceanMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

type DigitalOceanClusterProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

type DigitalOceanClusterProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	APIServerLB *godo.LoadBalancer `json:"apiServerLb,omitempty"`
}
