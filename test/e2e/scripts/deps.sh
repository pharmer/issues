#!/usr/bin/env bash

set -eoux pipefail

echo -e "\e[92mInstalling Dependencies ...\e[0m"

cd $HOME/go/src/github.com/pharmer/pharmer/
./hack/builddeps.sh
./hack/make.py
