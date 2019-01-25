package options

import "github.com/spf13/pflag"

type ClusterCreateOperation struct {
	OperationId string `json:"operation_id"`
}

func NewClusterCreateOperation() *ClusterCreateOperation  {
	return &ClusterCreateOperation{}
}

func (c *ClusterCreateOperation) AddFlags(fs *pflag.FlagSet)  {
	fs.StringVar(&c.OperationId, "operation-id", c.OperationId, "Operation id")
}
