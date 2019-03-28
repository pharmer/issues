#!/usr/bin/env bash

set -o xtrace
export KUBERNETES_VERSION="$1"
export CLUSTER_NAME="$2"
export PROVIDER="$3"
export ZONE="$4"
export NODE_TYPE="$5"
export REPLICAS="$6"
export CREDENTIAL=${PROVIDER}"_test"

echo -e "\e[92mCreating Cluster named $CLUSTER_NAME ...\e[0m"
pharmer create cluster ${CLUSTER_NAME} --provider=${PROVIDER} --zone=${ZONE} --nodes=${NODE_TYPE}=${REPLICAS} --credential-uid=${CREDENTIAL} --kubernetes-version=${KUBERNETES_VERSION}
pharmer apply ${CLUSTER_NAME} --v=10 || exit $?

echo -e "\e[92mSetting current kubernetes context to cluster $CLUSTER_NAME ...\e[0m"
pharmer use cluster ${CLUSTER_NAME}
