#!/usr/bin/env bash

set -o xtrace

export PROVIDER="$1"
export CREDENTIAL="$2""_test"
export FILE="$3"

echo -e "\e[92mCreating credential for $PROVIDER ...\e[0m"
if [[ "$FILE" == "" ]]; then
    pharmer create credential ${CREDENTIAL} -p ${PROVIDER} -l
else
    pharmer create credential ${CREDENTIAL} -p ${PROVIDER} --from-file ${FILE}
fi