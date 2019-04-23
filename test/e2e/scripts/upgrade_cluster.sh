#!/usr/bin/env bash

set -eoux pipefail

export KUBERNETES_VERSION="$1"
export CLUSTER_NAME="$2"

echo -e "\e[92mUpgrading Cluster $CLUSTER_NAME to $KUBERNETES_VERSION ...\e[0m"

pharmer edit cluster ${CLUSTER_NAME} --kubernetes-version=${KUBERNETES_VERSION}
pharmer apply ${CLUSTER_NAME}
