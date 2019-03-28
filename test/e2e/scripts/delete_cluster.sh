#!/usr/bin/env bash

set -o xtrace

export CLUSTER_NAME="$1"

echo -e "\e[92mDeleting Cluster $CLUSTER_NAME ...\e[0m"
pharmer delete cluster ${CLUSTER_NAME}
pharmer apply ${CLUSTER_NAME}
