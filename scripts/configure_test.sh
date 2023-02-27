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

# shellcheck source=./scripts/_assertions.sh
source _assertions.sh
# shellcheck source=scripts/_common.sh
source _common.sh

info "Assert WAN emulator image creation"
assert_non_empty "$(sudo docker images --filter reference=wanem --quiet)" "There is no WAN emulator Docker image created"

info "Assert KinD clusters creation"
assert_non_empty "$(sudo docker ps --filter label=io.x-k8s.kind.role=control-plane --quiet)" "There are no KinD clusters running"
