package xorm

import (
	"fmt"

	"github.com/go-xorm/xorm"
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/store"
)

type operationXormStore struct {
	engine *xorm.Engine
	//owner  string
}

var _ store.OperationStore = &operationXormStore{}

func (o *operationXormStore) Get(id string) (*api.Operation, error) {
	op := &api.Operation{
		Code: id,
	}
	fmt.Println("Xxxxxxxxxxx")
	has, err := o.engine.Get(op)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("not found")
	}
	return op, nil
}

func (o *operationXormStore) Update(obj *api.Operation) (*api.Operation, error) {
	o.engine.Update(obj)
	return obj, nil
}
