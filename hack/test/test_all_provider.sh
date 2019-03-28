#!/usr/bin/env bash

cd ${GOPATH}/src/github.com/pharmer/pharmer/test/e2e

export CURRENT_VERSION="$1"
export DESIRED_VERSION="$2"
export FILE="$3"

provider=(gce linode digitalocean packet vultr aws azure)
zone=(us-central1-f us-east nyc1 ewr1 6 us-east-1b westus2)
nodes=(n1-standard-2 g6-standard-2 2gb baremetal_0 94 t2.medium Standard_D1_v2)
file=(${FILE})

for i in 0 1 2 3 4 5 6
do
   go test --current-version=${CURRENT_VERSION} --desired-version=${DESIRED_VERSION} --provider=${provider[i]} --zone=${zone[i]} --nodes=${nodes[i]} --from-file=${file[0]} -timeout 1h
done
