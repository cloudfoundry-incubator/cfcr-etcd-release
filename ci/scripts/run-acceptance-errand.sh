#!/bin/bash

set -eu
source "$(dirname "$0")/helpers/expose-jumpbox.sh"
expose_jumpbox "${PWD}" "${JUMPBOX_SSH_KEY}" "${JUMPBOX_URL}"
bosh run-errand acceptance
