package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/go-nats-streaming"
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/apiserver/options"
	. "github.com/pharmer/pharmer/cloud"
	"log"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"strconv"
	"github.com/golang/glog"
	"time"
)

func (a *Apiserver) CreateCluster() error {

	sub, err := a.natsConn.Subscribe("create-cluster",  func(msg *stan.Msg) {
		fmt.Println("HERE...........")
		fmt.Printf("seq = %d [redelivered = %v, acked = false]\n", msg.Sequence, msg.Redelivered)

		operation := options.NewClusterCreateOperation()
		err := json.Unmarshal(msg.Data, &operation)
		if err != nil {
			glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
			return
		}
		if operation.OperationId == "" {
			// return error
		}


		obj, err := Store(a.ctx).Operations().Get(operation.OperationId)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(obj, obj.ID)
		fmt.Println(obj.Code, "XXXXXXXXXXXX", obj.State)

		if obj.State == api.OperationPending {
			obj.State = api.OperationRunning
			obj, err = Store(a.ctx).Operations().Update(obj)
			fmt.Println(obj)
			cluster, err := Store(a.ctx).Clusters().Get(strconv.Itoa(int(obj.ClusterID)))
			cluster.Spec.ClusterAPI = &clusterapi.Cluster{}

			cluster, err = Create(a.ctx, cluster, strconv.Itoa(int(obj.UserID)))
			if err != nil {
				fmt.Println(err)
				//term.Fatalln(err)
			}

			/*go func(o *opts.ApplyConfig) {
				acts, err := Apply(a.ctx, o)
				fmt.Println(acts, err)
			}(&opts.ApplyConfig{
				ClusterName: cluster.Name,
				Owner:       obj.UserID,
				DryRun:      false,
			})*/
		}
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ACK msg: %d", msg.Sequence)
		}

	}, stan.SetManualAckMode(), stan.AckWait(time.Second))
	if err != nil {
		return err
	}
	fmt.Println(sub)

	//defer LogCloser(sub)

	return nil
}
