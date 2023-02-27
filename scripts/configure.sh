#!/bin/bash
# SPDX-license-identifier: Apache-2.0
##############################################################################
# Copyright (c) 2023
# All rights reserved. This program and the accompanying materials
# are made available under the terms of the Apache License, Version 2.0
# which accompanies this distribution, and is available at
# http://www.apache.org/licenses/LICENSE-2.0
##############################################################################

set -o pipefail
set -o errexit
set -o nounset
[[ ${DEBUG:-false} != "true" ]] || set -o xtrace

# shellcheck source=scripts/_common.sh
source _common.sh
# shellcheck source=./scripts/_utils.sh
source _utils.sh

trap get_status ERR

wanem_img_name="wanem:0.0.1"

# Multi-cluster configuration
[[ -n "$(sudo docker images "$wanem_img_name" -q)" ]] || sudo docker build -t "$wanem_img_name" .
if ! sudo docker ps --format "{{.Image}}" | grep -q "kindest/node"; then
    # shellcheck disable=SC1091
    [ -f /etc/profile.d/path.sh ] && source /etc/profile.d/path.sh
    sudo -E "$(command -v go)" run ../... create --config ./config.yml --name test --wanem "$wanem_img_name"
    mkdir -p "$HOME/.kube"
    sudo chown -R "$USER" "$HOME/.kube/"
fi
