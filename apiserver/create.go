package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/nats-io/go-nats-streaming"
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/apiserver/options"
	. "github.com/pharmer/pharmer/cloud"
	opts "github.com/pharmer/pharmer/cloud/cmds/options"
	"github.com/pharmer/pharmer/notification"
	"strconv"
	"time"
)

func (a *Apiserver) CreateCluster() error {

	sub, err := a.natsConn.QueueSubscribe("create-cluster", "cluster-api-workers", func(msg *stan.Msg) {
		fmt.Printf("seq = %d [redelivered = %v, acked = false]\n", msg.Sequence, msg.Redelivered)

		operation := options.NewClusterCreateOperation()
		err := json.Unmarshal(msg.Data, &operation)
		if err != nil {
			glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
			return
		}
		if operation.OperationId == "" {
			err := fmt.Errorf("Operation id not  found")
			glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
			return
		}


		obj, err := Store(a.ctx).Operations().Get(operation.OperationId)
		if err != nil {
			glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
		}

		if obj.State == api.OperationPending {
			obj.State = api.OperationRunning
			obj, err = Store(a.ctx).Operations().Update(obj)
			if err != nil {
				glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
			}

			cluster, err := Store(a.ctx).Clusters().Get(strconv.Itoa(int(obj.ClusterID)))
			if err != nil {
				glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
			}

			cluster.InitClusterApi()


			noti := notification.NewNotifier(a.ctx, a.natsConn, strconv.Itoa(int(obj.ClusterID)))
			newCtx := WithLogger(a.ctx, noti)


			cluster, err = Create(newCtx, cluster, strconv.Itoa(int(obj.UserID)))
			if err != nil {
				glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
			}

			go func(o *opts.ApplyConfig, obj *api.Operation) {
				_, err := Apply(newCtx, o)
				if err != nil {
					glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
				}
				obj.State = api.OperationDone
				obj, err = Store(newCtx).Operations().Update(obj)
				if err != nil {
					glog.Errorf("seq = %d [redelivered = %v, data = %v, err = %v]\n", msg.Sequence, msg.Redelivered, msg.Data, err)
				}

				if err := msg.Ack(); err != nil {
					glog.Errorf("failed to ACK msg: %d", msg.Sequence)
				}

			}(&opts.ApplyConfig{
				ClusterName: cluster.Name, //strconv.Itoa(int(obj.ClusterID)),
				Owner:       strconv.Itoa(int(obj.UserID)),
				DryRun:      false,
			}, obj)
		}


	}, stan.SetManualAckMode(), stan.AckWait(time.Second))
	if err != nil {
		return err
	}
	fmt.Println(sub)

	//defer LogCloser(sub)

	return nil
}
