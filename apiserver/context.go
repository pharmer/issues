package apiserver

import "context"

type Apiserver struct {
	ctx context.Context
}

func New(ctx context.Context) *Apiserver  {
	return &Apiserver{ctx:ctx}
}



