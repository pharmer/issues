#!/usr/bin/env bash

set -eoux pipefail

export PROVIDER="$1"
export export CREDENTIAL="$2"
CREDENTIAL+="_test"

echo -e "\e[92mDeleting Credential $CREDENTIAL ...\e[0m"
pharmer delete cred $CREDENTIAL
