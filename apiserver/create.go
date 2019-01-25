package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/pharmer/pharmer/apiserver/options"
	"net/http"
)

func (a *Apiserver) CreateCluster(w http.ResponseWriter, r *http.Request) {
	operation := options.NewClusterCreateOperation()
	err := json.NewDecoder(r.Body).Decode(operation)
	if err != nil {
		// response invalid
	}
	fmt.Println(operation)


	
}
