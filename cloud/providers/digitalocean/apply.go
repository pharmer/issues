package digitalocean

import (
	"encoding/json"

	semver "github.com/appscode/go-version"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/clusterclient"

	//"context"
	"fmt"

	api "github.com/pharmer/pharmer/apis/v1beta1"
	. "github.com/pharmer/pharmer/cloud"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (cm *ClusterManager) Apply(in *api.Cluster, dryRun bool) ([]api.Action, error) {
	var err error
	var acts []api.Action

	if in.Status.Phase == "" {
		return nil, errors.Errorf("cluster `%s` is in unknown phase", cm.cluster.Name)
	}
	if in.Status.Phase == api.ClusterDeleted {
		return nil, nil
	}
	cm.cluster = in

	if cm.conn, err = PrepareCloud(cm.ctx, in.Name, cm.owner); err != nil {
		return nil, err
	}

	/*if err = cm.InitializeActuator(nil); err != nil {
		return nil, err
	}*/
	if cm.cluster.Status.Phase == api.ClusterUpgrading {
		return nil, errors.Errorf("cluster `%s` is upgrading. Retry after cluster returns to Ready state", cm.cluster.Name)
	}
	if cm.cluster.Status.Phase == api.ClusterReady {
		var kc kubernetes.Interface
		kc, err = cm.GetAdminClient()
		if err != nil {
			return nil, err
		}
		if upgrade, err := NewKubeVersionGetter(kc, cm.cluster).IsUpgradeRequested(); err != nil {
			return nil, err
		} else if upgrade {
			cm.cluster.Status.Phase = api.ClusterUpgrading
			Store(cm.ctx).Owner(cm.owner).Clusters().UpdateStatus(cm.cluster)
			return cm.applyUpgrade(dryRun)
		}
	}

	if cm.cluster.Status.Phase == api.ClusterPending {
		a, err := cm.applyCreate(dryRun)
		if err != nil {
			return nil, err
		}
		acts = append(acts, a...)
	}

	if cm.cluster.DeletionTimestamp != nil && cm.cluster.Status.Phase != api.ClusterDeleted {
		nodeGroups, err := Store(cm.ctx).Owner(cm.owner).MachineSet(cm.cluster.Name).List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		count := int32(0)
		for _, ng := range nodeGroups {
			ng.Spec.Replicas = &count
			_, err := Store(cm.ctx).Owner(cm.owner).MachineSet(cm.cluster.Name).Update(ng)
			if err != nil {
				return nil, err
			}
		}
	}

	{
		a, err := cm.applyScale(dryRun)
		if err != nil {
			return nil, err
		}
		acts = append(acts, a...)
	}

	if cm.cluster.DeletionTimestamp != nil && cm.cluster.Status.Phase != api.ClusterDeleted {
		a, err := cm.applyDelete(dryRun)
		if err != nil {
			return nil, err
		}
		acts = append(acts, a...)
	}
	return acts, nil
}

// Creates network, and creates ready master(s)
func (cm *ClusterManager) applyCreate(dryRun bool) (acts []api.Action, err error) {
	var found bool
	found, _, err = cm.conn.getPublicKey()
	if err != nil {
		return
	}
	if !found {
		acts = append(acts, api.Action{
			Action:   api.ActionAdd,
			Resource: "PublicKey",
			Message:  "Public key will be imported",
		})
		if !dryRun {
			cm.cluster.Status.Cloud.SShKeyExternalID, err = cm.conn.importPublicKey()
			if err != nil {
				return
			}
		}
	} else {
		acts = append(acts, api.Action{
			Action:   api.ActionNOP,
			Resource: "PublicKey",
			Message:  "Public key found",
		})
	}

	// ignore errors, since tags are simply informational.
	found, err = cm.conn.getTags()
	if err != nil {
		return
	}
	if !found {
		acts = append(acts, api.Action{
			Action:   api.ActionAdd,
			Resource: "Tag",
			Message:  fmt.Sprintf("Tag %s will be added", "KubernetesCluster:"+cm.cluster.Name),
		})
		if !dryRun {
			cm.conn.createTags()
		}
	} else {
		acts = append(acts, api.Action{
			Action:   api.ActionNOP,
			Resource: "Tag",
			Message:  fmt.Sprintf("Tag %s found", "KubernetesCluster:"+cm.cluster.Name),
		})
	}

	// -------------------------------------------------------------------ASSETS
	var machines []*clusterv1.Machine
	machines, err = Store(cm.ctx).Owner(cm.owner).Machine(cm.cluster.Name).List(metav1.ListOptions{})
	if err != nil {
		return
	}

	masterMachine, err := api.GetMasterMachine(machines)
	if err != nil {
		return
	}

	if d, _ := cm.conn.instanceIfExists(masterMachine); d == nil {
		Logger(cm.ctx).Info("Creating master instance")
		acts = append(acts, api.Action{
			Action:   api.ActionAdd,
			Resource: "MasterInstance",
			Message:  fmt.Sprintf("Master instance %s will be created", masterMachine.Name),
		})
		if !dryRun {
			var masterServer *api.NodeInfo
			nodeAddresses := make([]core.NodeAddress, 0)

			masterServer, err = cm.conn.CreateInstance(cm.cluster, masterMachine, "")
			if err != nil {
				return
			}

			if masterServer.PrivateIP != "" {
				nodeAddresses = append(nodeAddresses, core.NodeAddress{
					Type:    core.NodeInternalIP,
					Address: masterServer.PrivateIP,
				})
			}

			if masterServer.PublicIP != "" {
				nodeAddresses = append(nodeAddresses, core.NodeAddress{
					Type:    core.NodeExternalIP,
					Address: masterServer.PublicIP,
				})
			}
			if err = cm.cluster.SetClusterApiEndpoints(nodeAddresses); err != nil {
				return
			}

		}
	} else {
		acts = append(acts, api.Action{
			Action:   api.ActionNOP,
			Resource: "MasterInstance",
			Message:  fmt.Sprintf("master instance %v already exist", masterMachine.Name),
		})
	}

	if cm.cluster, err = Store(cm.ctx).Clusters().Update(cm.cluster); err != nil {
		return
	}

	var kc kubernetes.Interface
	kc, err = cm.GetAdminClient()
	if err != nil {
		return
	}
	// wait for nodes to start
	if err = WaitForReadyMaster(cm.ctx, kc); err != nil {
		return
	}

	// needed to get master_internal_ip
	cm.cluster.Status.Phase = api.ClusterReady
	if _, err = Store(cm.ctx).Clusters().UpdateStatus(cm.cluster); err != nil {
		return
	}

	// need to run ccm
	if err = CreateCredentialSecret(cm.ctx, kc, cm.cluster, cm.owner); err != nil {
		return
	}

	ca, err := NewClusterApi(cm.ctx, cm.cluster, "cloud-provider-system", kc)
	if err != nil {
		return acts, err
	}
	if err := ca.Apply(); err != nil {
		return acts, err
	}

	return acts, err
}

// Scales up/down regular node groups
func (cm *ClusterManager) applyScale(dryRun bool) (acts []api.Action, err error) {
	Logger(cm.ctx).Infoln("scaling machine set")

	var machineSets []*clusterv1.MachineSet
	var existingMachineSet []*clusterv1.MachineSet
	//var msc *clusterv1.MachineSet
	machineSets, err = Store(cm.ctx).Owner(cm.owner).MachineSet(cm.cluster.Name).List(metav1.ListOptions{})
	if err != nil {
		return
	}
	var bc clusterclient.Client
	bc, err = GetBooststrapClient(cm.ctx, cm.cluster, cm.owner)
	if err != nil {
		return nil, err
	}
	var data []byte
	for _, machineSet := range machineSets {
		if machineSet.DeletionTimestamp != nil {
			machineSet.DeletionTimestamp = nil
			if data, err = json.Marshal(machineSet); err != nil {
				return
			}

			if err = bc.Delete(string(data)); err != nil {
				return
			}
			if err = Store(cm.ctx).Owner(cm.owner).MachineSet(cm.cluster.Name).Delete(machineSet.Name); err != nil {
				return
			}
		}

		if existingMachineSet, err = bc.GetMachineSetObjects(); err != nil {
			return
		}
		if data, err = json.Marshal(machineSet); err != nil {
			return
		}
		found := false
		for _, ems := range existingMachineSet {
			if ems.Name == machineSet.Name {
				found = true
				if err = bc.Apply(string(data)); err != nil {
					return
				}
				break
			}
		}

		if !found {
			if err = bc.CreateMachineSetObjects([]*clusterv1.MachineSet{machineSet}, bc.GetContextNamespace()); err != nil {
				return
			}
		}
	}

	return
}

// Deletes master(s) and releases other cloud resources
func (cm *ClusterManager) applyDelete(dryRun bool) (acts []api.Action, err error) {
	Logger(cm.ctx).Infoln("deleting cluster")
	var found bool

	if cm.cluster.Status.Phase == api.ClusterReady {
		cm.cluster.Status.Phase = api.ClusterDeleting
	}
	_, err = Store(cm.ctx).Owner(cm.owner).Clusters().UpdateStatus(cm.cluster)
	if err != nil {
		return
	}

	var kc kubernetes.Interface
	kc, err = cm.GetAdminClient()
	if err != nil {
		return
	}
	var masterInstances *core.NodeList
	masterInstances, err = kc.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			api.RoleMasterKey: "",
		}).String(),
	})
	if err != nil && !kerr.IsNotFound(err) {
		Logger(cm.ctx).Infof("master instance not found. Reason: %v", err)
	} else if err == nil {
		acts = append(acts, api.Action{
			Action:   api.ActionDelete,
			Resource: "MasterInstance",
			Message:  fmt.Sprintf("Will delete master instance with name %v-master", cm.cluster.Name),
		})
		if !dryRun {
			for _, mi := range masterInstances.Items {
				err = cm.conn.DeleteInstanceByProviderID(mi.Spec.ProviderID)
				if err != nil {
					Logger(cm.ctx).Infof("Failed to delete instance %s. Reason: %s", mi.Spec.ProviderID, err)
				}
			}
			/*if masterNG.Spec.Template.Spec.ExternalIPType == api.IPTypeReserved {
				for _, addr := range masterInstance.Status.Addresses {
					if addr.Type == core.NodeExternalIP {
						err = cm.conn.releaseReservedIP(addr.Address)
						if err != nil {
							return
						}
					}
				}
			}*/
		}
	}

	// delete by tag
	tag := "KubernetesCluster:" + cm.cluster.Name
	_, err = cm.conn.client.Droplets.DeleteByTag(cm.ctx, tag)
	if err != nil {
		Logger(cm.ctx).Infof("Failed to delete resources by tag %s. Reason: %s", tag, err)
	}
	Logger(cm.ctx).Infof("Deleted droplet by tag %s", tag)

	// Delete SSH key
	found, _, err = cm.conn.getPublicKey()
	if err != nil {
		return
	}
	if found {
		acts = append(acts, api.Action{
			Action:   api.ActionAdd,
			Resource: "PublicKey",
			Message:  "Public key will be deleted",
		})
		if !dryRun {
			err = cm.conn.deleteSSHKey()
			if err != nil {
				return
			}
		}
	} else {
		acts = append(acts, api.Action{
			Action:   api.ActionNOP,
			Resource: "PublicKey",
			Message:  "Public key not found",
		})
	}
	/*if IsHASetup(cm.cluster) {
		cm.conn.deleteLoadBalancer(cm.ctx, cm.namer.LoadBalancerName())
	}
	*/
	// Failed
	cm.cluster.Status.Phase = api.ClusterDeleted
	_, err = Store(cm.ctx).Owner(cm.owner).Clusters().UpdateStatus(cm.cluster)
	if err != nil {
		return
	}

	Logger(cm.ctx).Infof("Cluster %v deletion is deleted successfully", cm.cluster.Name)
	return
}

func (cm *ClusterManager) applyUpgrade(dryRun bool) (acts []api.Action, err error) {
	var kc kubernetes.Interface
	if kc, err = cm.GetAdminClient(); err != nil {
		return
	}

	var masterMachine *clusterv1.Machine
	masterName := fmt.Sprintf("%v-master", cm.cluster.Name)
	masterMachine, err = Store(cm.ctx).Machine(cm.cluster.Name).Get(masterName)
	if err != nil {
		return nil, err
	}

	masterMachine.Spec.Versions.ControlPlane = cm.cluster.Spec.Config.KubernetesVersion
	masterMachine.Spec.Versions.Kubelet = cm.cluster.Spec.Config.KubernetesVersion

	var bc clusterclient.Client
	bc, err = GetBooststrapClient(cm.ctx, cm.cluster, cm.owner)
	if err != nil {
		return nil, err
	}

	var data []byte
	if data, err = json.Marshal(masterMachine); err != nil {
		return
	}
	if err = bc.Apply(string(data)); err != nil {
		return
	}

	// Wait until master updated
	desiredVersion, _ := semver.NewVersion(cm.cluster.ClusterConfig().KubernetesVersion)
	if err = WaitForReadyMasterVersion(cm.ctx, kc, desiredVersion); err != nil {
		return
	}
	// wait for nodes to start
	if err = WaitForReadyMaster(cm.ctx, kc); err != nil {
		return
	}

	var machineSets []*clusterv1.MachineSet
	machineSets, err = Store(cm.ctx).MachineSet(cm.cluster.Name).List(metav1.ListOptions{})
	if err != nil {
		return
	}

	for _, machineSet := range machineSets {
		machineSet.Spec.Template.Spec.Versions.Kubelet = cm.cluster.Spec.Config.KubernetesVersion
		if data, err = json.Marshal(machineSet); err != nil {
			return
		}
		if err = bc.Apply(string(data)); err != nil {
			return
		}
	}

	if !dryRun {
		cm.cluster.Status.Phase = api.ClusterReady
		if _, err = Store(cm.ctx).Owner(cm.owner).Clusters().UpdateStatus(cm.cluster); err != nil {
			return
		}
	}
	return
}
