package v1beta1

import (
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	DigitalOceanProviderGroupName  = "digitaloceanproviderconfig"
	DigitalOceanProviderKind       = "DigitaloceanClusterProviderConfig"
	DigitalOceanProviderApiVersion = "v1alpha1"

	GKEProviderGroupName  = "gkeproviderconfig"
	GKEProviderKind       = "GKEProviderConfig"
	GKEProviderApiVersion = "v1alpha1"
)

// DigitalOceanMachineProviderConfig contains Config for DigitalOcean machines.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DigitalOceanMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

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

func (c *Cluster) DigitalOceanProviderConfig(cluster *clusterapi.Cluster) *DigitalOceanMachineProviderConfig {
	//providerConfig providerConfig
	raw := cluster.Spec.ProviderSpec.Value.Raw
	providerConfig := &DigitalOceanMachineProviderConfig{}
	err := json.Unmarshal(raw, providerConfig)
	if err != nil {
		fmt.Println("Unable to unmarshal provider config: %v", err)
	}
	return providerConfig
}

func (c *Cluster) SetDigitalOceanProviderConfig(cluster *clusterapi.Cluster, config *ClusterConfig) error {
	conf := &DigitalOceanMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: DigitalOceanProviderGroupName + "/" + DigitalOceanProviderApiVersion,
			Kind:       DigitalOceanProviderKind,
		},
	}
	bytes, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("Unable to marshal provider config: %v", err)
		return err
	}
	cluster.Spec.ProviderSpec = clusterapi.ProviderSpec{
		Value: &runtime.RawExtension{
			Raw: bytes,
		},
	}
	return nil
}

type GKEMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Roles []MachineRole `json:"roles,omitempty"`

	Zone        string `json:"zone"`
	MachineType string `json:"machineType"`

	// The name of the OS to be installed on the machine.
	OS    string `json:"os,omitempty"`
	Disks []Disk `json:"disks,omitempty"`
}

type Disk struct {
	InitializeParams DiskInitializeParams `json:"initializeParams"`
}

type DiskInitializeParams struct {
	DiskSizeGb int64  `json:"diskSizeGb"`
	DiskType   string `json:"diskType"`
}

func (c *Cluster) GKEProviderConfig(raw []byte) *GKEMachineProviderSpec {
	//providerConfig providerConfig
	providerConfig := &GKEMachineProviderSpec{}
	err := json.Unmarshal(raw, providerConfig)
	if err != nil {
		fmt.Println("Unable to unmarshal provider config: %v", err)
	}
	return providerConfig
}

func (c *Cluster) SetGKEProviderConfig(cluster *clusterapi.Cluster, config *ClusterConfig) error {
	conf := &GKEMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GKEProviderGroupName + "/" + GKEProviderApiVersion,
			Kind:       GKEProviderKind,
		},
	}
	bytes, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("Unable to marshal provider config: %v", err)
		return err
	}
	cluster.Spec.ProviderSpec = clusterapi.ProviderSpec{
		Value: &runtime.RawExtension{
			Raw: bytes,
		},
	}
	return nil
}
