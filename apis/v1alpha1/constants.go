package v1alpha1

import "time"

const (
	RoleMaster    = "master"
	RoleNode      = "node"
	RoleKeyPrefix = "node-role.kubernetes.io/"
	RoleMasterKey = RoleKeyPrefix + RoleMaster
	RoleNodeKey   = RoleKeyPrefix + RoleNode

	KubeadmVersionKey = "cluster.pharmer.io/kubeadm-version"
	NodePoolKey       = "cluster.pharmer.io/pool"
	KubeSystem_App    = "k8s-app"

	HostnameKey     = "kubernetes.io/hostname"
	ArchKey         = "beta.kubernetes.io/arch"
	InstanceTypeKey = "beta.kubernetes.io/instance-type"
	OSKey           = "beta.kubernetes.io/os"
	RegionKey       = "failure-domain.beta.kubernetes.io/region"
	ZoneKey         = "failure-domain.beta.kubernetes.io/zone"

	// CoreDNS defines a variable used internally when referring to the CoreDNS addon for a cluster
	CoreDNS = "coredns"

	TokenDuration_10yr = 10 * 365 * 24 * time.Hour

	// ref: https://github.com/kubernetes/kubeadm/issues/629
	DeprecatedV19AdmissionControl = "NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,ValidatingAdmissionWebhook,DefaultTolerationSeconds,MutatingAdmissionWebhook,ResourceQuota"
	DefaultV19AdmissionControl    = "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ValidatingAdmissionWebhook,DefaultTolerationSeconds,MutatingAdmissionWebhook,ResourceQuota"
	DefaultV111AdmissionControl   = "Initializers,NodeRestriction,NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ValidatingAdmissionWebhook,DefaultTolerationSeconds,MutatingAdmissionWebhook,ResourceQuota"
)
