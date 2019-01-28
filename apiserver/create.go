package apiserver

import (
	"encoding/json"
	"fmt"
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/apiserver/options"
	. "github.com/pharmer/pharmer/cloud"
	"net/http"
	"strconv"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (a *Apiserver) CreateCluster(w http.ResponseWriter, r *http.Request) {
	operation := options.NewClusterCreateOperation()
	err := json.NewDecoder(r.Body).Decode(operation)
	if err != nil {
		// response invalid
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
		fmt.Println(cluster, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
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



	
}
