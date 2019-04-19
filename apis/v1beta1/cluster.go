package v1beta1

import (
	"fmt"

	"github.com/appscode/go-version"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	ResourceCodeCluster = ""
	ResourceKindCluster = "Cluster"
	ResourceNameCluster = "cluster"
	ResourceTypeCluster = "clusters"

	DefaultKubernetesBindPort = 6443
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Cluster struct {
	metav1.TypeMeta   `json:",inline,omitempty,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              PharmerClusterSpec   `json:"spec,omitempty,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            PharmerClusterStatus `json:"status,omitempty,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type PharmerClusterSpec struct {
	ClusterAPI *clusterapi.Cluster `json:"clusterApi,omitempty" protobuf:"bytes,1,opt,name=clusterApi"`
	Config     *ClusterConfig      `json:"config,omitempty" protobuf:"bytes,2,opt,name=config"`
}

type ClusterConfig struct {
	MasterCount          int       `json:"masterCount"`
	Cloud                CloudSpec `json:"cloud" protobuf:"bytes,1,opt,name=cloud"`
	KubernetesVersion    string    `json:"kubernetesVersion,omitempty" protobuf:"bytes,4,opt,name=kubernetesVersion"`
	Locked               bool      `json:"locked,omitempty" protobuf:"varint,5,opt,name=locked"`
	CACertName           string    `json:"caCertName,omitempty" protobuf:"bytes,6,opt,name=caCertName"`
	FrontProxyCACertName string    `json:"frontProxyCACertName,omitempty" protobuf:"bytes,7,opt,name=frontProxyCACertName"`
	CredentialName       string    `json:"credentialName,omitempty" protobuf:"bytes,8,opt,name=credentialName"`

	KubeletExtraArgs           map[string]string `json:"kubeletExtraArgs,omitempty" protobuf:"bytes,9,rep,name=kubeletExtraArgs"`
	APIServerExtraArgs         map[string]string `json:"apiServerExtraArgs,omitempty" protobuf:"bytes,10,rep,name=apiServerExtraArgs"`
	ControllerManagerExtraArgs map[string]string `json:"controllerManagerExtraArgs,omitempty" protobuf:"bytes,11,rep,name=controllerManagerExtraArgs"`
	SchedulerExtraArgs         map[string]string `json:"schedulerExtraArgs,omitempty" protobuf:"bytes,12,rep,name=schedulerExtraArgs"`
	AuthorizationModes         []string          `json:"authorizationModes,omitempty" protobuf:"bytes,13,rep,name=authorizationModes"`
	APIServerCertSANs          []string          `json:"apiServerCertSANs,omitempty" protobuf:"bytes,14,rep,name=apiServerCertSANs"`
}

type API struct {
	// AdvertiseAddress sets the address for the API server to advertise.
	AdvertiseAddress string `json:"advertiseAddress" protobuf:"bytes,1,opt,name=advertiseAddress"`
	// BindPort sets the secure port for the API Server to bind to
	BindPort int32 `json:"bindPort" protobuf:"varint,2,opt,name=bindPort"`
}

type CloudSpec struct {
	CloudProvider        string      `json:"cloudProvider,omitempty" protobuf:"bytes,1,opt,name=cloudProvider"`
	Project              string      `json:"project,omitempty" protobuf:"bytes,2,opt,name=project"`
	Region               string      `json:"region,omitempty" protobuf:"bytes,3,opt,name=region"`
	Zone                 string      `json:"zone,omitempty" protobuf:"bytes,4,opt,name=zone"` // master needs it for ossec
	InstanceImage        string      `json:"instanceImage,omitempty" protobuf:"bytes,5,opt,name=instanceImage"`
	OS                   string      `json:"os,omitempty" protobuf:"bytes,6,opt,name=os"`
	InstanceImageProject string      `json:"instanceImageProject,omitempty" protobuf:"bytes,7,opt,name=instanceImageProject"`
	NetworkProvider      string      `json:"networkProvider,omitempty" protobuf:"bytes,8,opt,name=networkProvider"` // kubenet, flannel, calico, opencontrail
	CCMCredentialName    string      `json:"ccmCredentialName,omitempty" protobuf:"bytes,9,opt,name=ccmCredentialName"`
	SSHKeyName           string      `json:"sshKeyName,omitempty" protobuf:"bytes,10,opt,name=sshKeyName"`
	AWS                  *AWSSpec    `json:"aws,omitempty" protobuf:"bytes,11,opt,name=aws"`
	GCE                  *GoogleSpec `json:"gce,omitempty" protobuf:"bytes,12,opt,name=gce"`
	Azure                *AzureSpec  `json:"azure,omitempty" protobuf:"bytes,13,opt,name=azure"`
	Linode               *LinodeSpec `json:"linode,omitempty" protobuf:"bytes,14,opt,name=linode"`
	GKE                  *GKESpec    `json:"gke,omitempty" protobuf:"bytes,15,opt,name=gke"`
	//DigitalOcean         *DigitalOceanMachineProviderConfig `json:"digitalocean,omitempty" protobuf:"bytes,16,opt,name=digitalocean"`
	Dokube *DokubeSpec `json:"dokube,omitempty" protobuf:"bytes,15,opt,name=dokube"`
}

type AWSSpec struct {
	// aws:TAG KubernetesCluster => clusterid
	IAMProfileMaster  string `json:"iamProfileMaster,omitempty" protobuf:"bytes,1,opt,name=iamProfileMaster"`
	IAMProfileNode    string `json:"iamProfileNode,omitempty" protobuf:"bytes,2,opt,name=iamProfileNode"`
	MasterSGName      string `json:"masterSGName,omitempty" protobuf:"bytes,3,opt,name=masterSGName"`
	NodeSGName        string `json:"nodeSGName,omitempty" protobuf:"bytes,4,opt,name=nodeSGName"`
	BastionSGName     string `json:"bastionSGName,omitempty"`
	VpcCIDR           string `json:"vpcCIDR,omitempty" protobuf:"bytes,5,opt,name=vpcCIDR"`
	VpcCIDRBase       string `json:"vpcCIDRBase,omitempty" protobuf:"bytes,6,opt,name=vpcCIDRBase"`
	MasterIPSuffix    string `json:"masterIPSuffix,omitempty" protobuf:"bytes,7,opt,name=masterIPSuffix"`
	PrivateSubnetCIDR string `json:"privateSubnetCidr,omitempty" protobuf:"bytes,8,opt,name=privateSubnetCidr"`
	PublicSubnetCIDR  string `json:"publicSubnetCidr,omitempty" protobuf:"bytes,9,opt,name=publicSubnetCidr"`
}

type GoogleSpec struct {
	NetworkName string   `gcfg:"network-name" ini:"network-name,omitempty" protobuf:"bytes,1,opt,name=networkName"`
	NodeTags    []string `gcfg:"node-tags" ini:"node-tags,omitempty,omitempty" protobuf:"bytes,2,rep,name=nodeTags"`
	// gce
	// NODE_SCOPES="${NODE_SCOPES:-compute-rw,monitoring,logging-write,storage-ro}"
	NodeScopes []string `json:"nodeScopes,omitempty" protobuf:"bytes,3,rep,name=nodeScopes"`
}

type GCECloudConfig struct {
	TokenURL           string   `gcfg:"token-url" ini:"token-url,omitempty" protobuf:"bytes,1,opt,name=tokenURL"`
	TokenBody          string   `gcfg:"token-body" ini:"token-body,omitempty" protobuf:"bytes,2,opt,name=tokenBody"`
	ProjectID          string   `gcfg:"project-id" ini:"project-id,omitempty" protobuf:"bytes,3,opt,name=projectID"`
	NetworkName        string   `gcfg:"network-name" ini:"network-name,omitempty" protobuf:"bytes,4,opt,name=networkName"`
	NodeTags           []string `gcfg:"node-tags" ini:"node-tags,omitempty,omitempty" protobuf:"bytes,5,rep,name=nodeTags"`
	NodeInstancePrefix string   `gcfg:"node-instance-prefix" ini:"node-instance-prefix,omitempty,omitempty" protobuf:"bytes,6,opt,name=nodeInstancePrefix"`
	Multizone          bool     `gcfg:"multizone" ini:"multizone,omitempty" protobuf:"varint,7,opt,name=multizone"`
}

type GKESpec struct {
	UserName    string `json:"userName,omitempty" protobuf:"bytes,1,opt,name=userName"`
	Password    string `json:"password,omitempty" protobuf:"bytes,2,opt,name=password"`
	NetworkName string `json:"networkName,omitempty" protobuf:"bytes,3,opt,name=networkName"`
}

type AzureSpec struct {
	InstanceImageVersion   string `json:"instanceImageVersion,omitempty" protobuf:"bytes,1,opt,name=instanceImageVersion"`
	RootPassword           string `json:"rootPassword,omitempty" protobuf:"bytes,2,opt,name=rootPassword"`
	VPCCIDR                string `json:"vpcCIDR"`
	ControlPlaneSubnetCIDR string `json:"controlPlaneSubnetCIDR"`
	NodeSubnetCIDR         string `json:"nodeSubnetCIDR"`
	InternalLBIPAddress    string `json:"internalLBIPAddress"`
	AzureDNSZone           string `json:"azureDNSZone"`
	SubnetCIDR             string `json:"subnetCidr,omitempty" protobuf:"bytes,3,opt,name=subnetCidr"`
	ResourceGroup          string `json:"resourceGroup,omitempty" protobuf:"bytes,4,opt,name=resourceGroup"`
	SubnetName             string `json:"subnetName,omitempty" protobuf:"bytes,5,opt,name=subnetName"`
	SecurityGroupName      string `json:"securityGroupName,omitempty" protobuf:"bytes,6,opt,name=securityGroupName"`
	VnetName               string `json:"vnetName,omitempty" protobuf:"bytes,7,opt,name=vnetName"`
	RouteTableName         string `json:"routeTableName,omitempty" protobuf:"bytes,8,opt,name=routeTableName"`
	StorageAccountName     string `json:"azureStorageAccountName,omitempty" protobuf:"bytes,9,opt,name=azureStorageAccountName"`
	SubscriptionID         string `json:"subscriptionID"`
}

// ref: https://github.com/kubernetes/kubernetes/blob/8b9f0ea5de2083589f3b9b289b90273556bc09c4/pkg/cloudprovider/providers/azure/azure.go#L56
type AzureCloudConfig struct {
	TenantID           string `json:"tenantId,omitempty" protobuf:"bytes,1,opt,name=tenantId"`
	SubscriptionID     string `json:"subscriptionId,omitempty" protobuf:"bytes,2,opt,name=subscriptionId"`
	AadClientID        string `json:"aadClientId,omitempty" protobuf:"bytes,3,opt,name=aadClientId"`
	AadClientSecret    string `json:"aadClientSecret,omitempty" protobuf:"bytes,4,opt,name=aadClientSecret"`
	ResourceGroup      string `json:"resourceGroup,omitempty" protobuf:"bytes,5,opt,name=resourceGroup"`
	Location           string `json:"location,omitempty" protobuf:"bytes,6,opt,name=location"`
	SubnetName         string `json:"subnetName,omitempty" protobuf:"bytes,7,opt,name=subnetName"`
	SecurityGroupName  string `json:"securityGroupName,omitempty" protobuf:"bytes,8,opt,name=securityGroupName"`
	VnetName           string `json:"vnetName,omitempty" protobuf:"bytes,9,opt,name=vnetName"`
	RouteTableName     string `json:"routeTableName,omitempty" protobuf:"bytes,10,opt,name=routeTableName"`
	StorageAccountName string `json:"storageAccountName,omitempty" protobuf:"bytes,11,opt,name=storageAccountName"`
}

type LinodeSpec struct {
	// Linode
	RootPassword string `json:"rootPassword,omitempty" protobuf:"bytes,1,opt,name=rootPassword"`
	KernelId     string `json:"kernelId,omitempty" protobuf:"varint,2,opt,name=kernelId"`
}

type LinodeCloudConfig struct {
	Token string `json:"token,omitempty" protobuf:"bytes,1,opt,name=token"`
	Zone  string `json:"zone,omitempty" protobuf:"bytes,2,opt,name=zone"`
}

type PacketCloudConfig struct {
	Project string `json:"project,omitempty" protobuf:"bytes,1,opt,name=project"`
	ApiKey  string `json:"apiKey,omitempty" protobuf:"bytes,2,opt,name=apiKey"`
	Zone    string `json:"zone,omitempty" protobuf:"bytes,3,opt,name=zone"`
}

type VultrCloudConfig struct {
	Token string `json:"token,omitempty" protobuf:"bytes,1,opt,name=token"`
}

type DokubeSpec struct {
	ClusterID string `json:"clusterID,omitempty" protobuf:"bytes,3,opt,name=clusterID"`
}

// ClusterPhase is a label for the condition of a Cluster at the current time.
type ClusterPhase string

// These are the valid statuses of Cluster.
const (
	ClusterPending   ClusterPhase = "Pending"
	ClusterReady     ClusterPhase = "Ready"
	ClusterDeleting  ClusterPhase = "Deleting"
	ClusterDeleted   ClusterPhase = "Deleted"
	ClusterUpgrading ClusterPhase = "Upgrading"
)

type CloudStatus struct {
	SShKeyExternalID string       `json:"sshKeyExternalID,omitempty" protobuf:"bytes,1,opt,name=sshKeyExternalID"`
	AWS              *AWSStatus   `json:"aws,omitempty" protobuf:"bytes,2,opt,name=aws"`
	EKS              *EKSStatus   `json:"eks,omitempty" protobuf:"bytes,2,opt,name=eks"`
	LoadBalancer     LoadBalancer `json:"loadBalancer,omitempty"`
}

type LoadBalancer struct {
	DNS  string `json:"dns"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type AWSStatus struct {
	MasterSGId  string `json:"masterSGID,omitempty" protobuf:"bytes,1,opt,name=masterSGID"`
	NodeSGId    string `json:"nodeSGID,omitempty" protobuf:"bytes,2,opt,name=nodeSGID"`
	BastionSGId string `json:"bastionSGID,omitempty"`

	// Depricaed
	// TODO: REMOVE
	LBDNS string `json:"lbDNS,omitempty"`
}

type EKSStatus struct {
	SecurityGroup string `json:"securityGroup,omitempty" protobuf:"bytes,1,opt,name=securityGroup"`
	VpcId         string `json:"vpcID,omitempty" protobuf:"bytes,2,opt,name=vpcID"`
	SubnetId      string `json:"subnetID,omitempty" protobuf:"bytes,3,opt,name=subnetID"`
	RoleArn       string `json:"roleArn,omitempty" protobuf:"bytes,4,opt,name=roleArn"`
}

type PharmerClusterStatus struct {
	Phase  ClusterPhase `json:"phase,omitempty,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=ClusterPhase"`
	Reason string       `json:"reason,omitempty,omitempty" protobuf:"bytes,2,opt,name=reason"`
	Cloud  CloudStatus  `json:"cloud,omitempty" protobuf:"bytes,4,opt,name=cloud"`
	//ReservedIPs  []ReservedIP       `json:"reservedIP,omitempty" protobuf:"bytes,6,rep,name=reservedIP"`
}

type ReservedIP struct {
	IP   string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=ip"`
	ID   string `json:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
}

func (c *Cluster) ClusterConfig() *ClusterConfig {
	return c.Spec.Config
}

func (c *Cluster) APIServerURL() string {
	for _, addr := range c.Spec.ClusterAPI.Status.APIEndpoints {
		if addr.Port == 0 {
			return fmt.Sprintf("https://%s", addr.Host)
		} else {
			return fmt.Sprintf("https://%s:%d", addr.Host, addr.Port)
		}

	}
	return ""
}

func (c *Cluster) SetClusterApiEndpoints(addresses []core.NodeAddress) error {
	m := map[core.NodeAddressType]string{}
	for _, addr := range addresses {
		m[addr.Type] = addr.Address

	}
	if u, found := m[core.NodeExternalIP]; found {
		c.Spec.ClusterAPI.Status.APIEndpoints = append(c.Spec.ClusterAPI.Status.APIEndpoints, clusterapi.APIEndpoint{
			Host: u,
			Port: int(DefaultKubernetesBindPort),
		})
		return nil
	}
	if u, found := m[core.NodeExternalDNS]; found {
		c.Spec.ClusterAPI.Status.APIEndpoints = append(c.Spec.ClusterAPI.Status.APIEndpoints, clusterapi.APIEndpoint{
			Host: u,
			Port: int(DefaultKubernetesBindPort),
		})
		return nil
	}
	return fmt.Errorf("No cluster api endpoint found")
}

func (c *Cluster) APIServerAddress() string {
	endpoints := c.Spec.ClusterAPI.Status.APIEndpoints
	if len(endpoints) == 0 {
		return ""
	}
	ep := endpoints[0]
	if ep.Port == 0 {
		return ep.Host
	} else {
		return fmt.Sprintf("%s:%d", ep.Host, ep.Port)
	}

}

func (c *Cluster) SetNetworkingDefaults(provider string) {
	clusterSpec := &c.Spec.ClusterAPI.Spec
	if len(clusterSpec.ClusterNetwork.Services.CIDRBlocks) == 0 {
		clusterSpec.ClusterNetwork.Services.CIDRBlocks = []string{kubeadmapi.DefaultServicesSubnet}
	}
	if clusterSpec.ClusterNetwork.ServiceDomain == "" {
		clusterSpec.ClusterNetwork.ServiceDomain = kubeadmapi.DefaultServiceDNSDomain
	}
	if len(clusterSpec.ClusterNetwork.Pods.CIDRBlocks) == 0 {
		// https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network
		podSubnet := ""
		switch provider {
		case PodNetworkCalico:
			podSubnet = "192.168.0.0/16"
		case PodNetworkFlannel:
			podSubnet = "10.244.0.0/16"
		case PodNetworkCanal:
			podSubnet = "10.244.0.0/16"
		}
		clusterSpec.ClusterNetwork.Pods.CIDRBlocks = []string{podSubnet}
	}
}

func (c *Cluster) InitClusterApi() {
	c.Spec.ClusterAPI = &clusterapi.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.Name,
		},
	}
}

func (c Cluster) IsMinorVersion(in string) bool {
	v, err := version.NewVersion(c.Spec.Config.KubernetesVersion)
	if err != nil {
		return false
	}
	minor := v.ToMutator().ResetMetadata().ResetPrerelease().ResetPatch().String()

	inVer, err := version.NewVersion(in)
	if err != nil {
		return false
	}
	return inVer.String() == minor
}

func (c Cluster) IsLessThanVersion(in string) bool {
	v, err := version.NewVersion(c.Spec.Config.KubernetesVersion)
	if err != nil {
		return false
	}
	inVer, err := version.NewVersion(in)
	if err != nil {
		return false
	}
	return v.LessThan(inVer)
}
