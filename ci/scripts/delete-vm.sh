#!/bin/bash

set -eu

source "$(dirname "$0")/helpers/expose-jumpbox.sh"
expose_jumpbox "${PWD}" "${JUMPBOX_SSH_KEY}" "${JUMPBOX_URL}"

vm_id="$(bosh vms | grep vm | tail -1 | awk '{print $5}')"
bosh -n delete-vm "${vm_id}"