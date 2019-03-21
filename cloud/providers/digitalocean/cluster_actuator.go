package digitalocean

import (
	"context"
	"reflect"

	"github.com/appscode/go/log"
	doCapi "github.com/pharmer/pharmer/apis/v1beta1/digitalocean"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/controller/cluster"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, func(ctx context.Context, m manager.Manager, owner string) error {
		actuator := NewClusterActuator(m, ClusterActuatorParams{
			Ctx:           ctx,
			EventRecorder: m.GetRecorder(Recorder),
			Scheme:        m.GetScheme(),
			Owner:         owner,
		})
		return cluster.AddWithActuator(m, actuator)
	})

}

type ClusterActuator struct {
	ctx           context.Context
	client        client.Client
	eventRecorder record.EventRecorder
	scheme        *runtime.Scheme
	owner         string
	conn          *cloudConnector
}

type ClusterActuatorParams struct {
	Ctx            context.Context
	EventRecorder  record.EventRecorder
	Scheme         *runtime.Scheme
	Owner          string
	CloudConnector *cloudConnector
}

func NewClusterActuator(m manager.Manager, params ClusterActuatorParams) *ClusterActuator {
	return &ClusterActuator{
		ctx:           params.Ctx,
		client:        m.GetClient(),
		eventRecorder: params.EventRecorder,
		scheme:        params.Scheme,
		owner:         params.Owner,
		conn:          params.CloudConnector,
	}
}

func (a *ClusterActuator) Reconcile(cluster *clusterapi.Cluster) error {
	log.Infoln("Reconciling cluster", cluster.Name)

	conn, err := PrepareCloud(a.ctx, cluster.Name, a.owner)
	if err != nil {
		log.Debugln("Error creating cloud connector", err)
		return err
	}
	a.conn = conn

	// TODO move to reconcileLoadBalance() func if more things are added here
	lb, err := a.conn.lbByName(context.Background(), a.conn.namer.LoadBalancerName())
	if err == errLBNotFound {
		lb, err = a.conn.createLoadBalancer(context.Background(), a.conn.namer.LoadBalancerName())
		if err != nil {
			log.Debugln("error creating load balancer", err)
			return err
		}
		log.Infof("created load balancer %q for cluster %q", a.conn.namer.LoadBalancerName(), cluster.Name)

		cluster.Status.APIEndpoints = []clusterapi.APIEndpoint{
			{
				Host: lb.IP,
				Port: v1beta1.DefaultAPIBindPort,
			},
		}
	} else if err != nil {
		log.Debugln("error finding load balancer", err)
		return err
	}

	// now check load balancer specs
	defaultSpecs, err := a.conn.buildLoadBalancerRequest(a.conn.namer.LoadBalancerName())
	if err != nil {
		log.Debugln("error getting default lb specs")
		return err
	}

	updateRequired := false

	if lb.Algorithm != defaultSpecs.Algorithm {
		updateRequired = true
	}
	if lb.Region.Slug != defaultSpecs.Region {
		updateRequired = true
	}
	if !reflect.DeepEqual(lb.ForwardingRules, defaultSpecs.ForwardingRules) {
		updateRequired = true
	}
	if !reflect.DeepEqual(lb.HealthCheck, defaultSpecs.HealthCheck) {
		updateRequired = true
	}
	if !reflect.DeepEqual(lb.StickySessions, defaultSpecs.StickySessions) {
		updateRequired = true
	}
	if lb.RedirectHttpToHttps != defaultSpecs.RedirectHttpToHttps {
		updateRequired = true
	}

	if updateRequired {
		log.Infoln("load balancer specs changed, updating lb")
		lb, _, err = a.conn.client.LoadBalancers.Update(context.Background(), lb.ID, defaultSpecs)
		if err != nil {
			log.Debugln("error updating load balancer", err)
			return err
		}
	}

	status, err := doCapi.ClusterStatusFromProviderStatus(cluster.Status.ProviderStatus)
	if err != nil {
		log.Debugln("Error getting provider status", err)
		return err
	}
	status.APIServerLB = lb

	log.Infoln("Reconciled cluster successfully")
	return nil
}

func (a *ClusterActuator) Delete(cluster *clusterapi.Cluster) error {
	log.Infoln("Delete cluster not implemented")

	return nil
}
