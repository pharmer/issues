package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/pharmer/pharmer/apiserver/options"
	. "github.com/pharmer/pharmer/cloud"
	"net/http"
	opts "github.com/pharmer/pharmer/cloud/cmds/options"
	api "github.com/pharmer/pharmer/apis/v1beta1"
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
		// error
	}

	if obj.State == api.OperationPending {
		obj.State = api.OperationProgress
		Store(a.ctx).Operations().Update(obj)
		cluster, err := Store(a.ctx).Clusters().Get(obj.ClusterID)

		cluster, err = Create(a.ctx, cluster, obj.UserID)
		if err != nil {
			//term.Fatalln(err)
		}

		go func(o *opts.ApplyConfig) {
			acts, err := Apply(a.ctx, o)
			fmt.Println(acts, err)
		}(&opts.ApplyConfig{
			ClusterName: cluster.Name,
			Owner:       obj.UserID,
			DryRun:      false,
		})
	}



	
}
