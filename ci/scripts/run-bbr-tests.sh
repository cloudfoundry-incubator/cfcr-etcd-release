#!/bin/bash

set -euo pipefail

tar xvf bbr-release/bbr-*.tar
mv releases/bbr /usr/local/bin/bbr

jumpbox_private_key="$(mktemp)"
BOSH_CA_CERT_PATH="$(mktemp)"

echo "${JUMPBOX_PRIVATE_KEY}" > "${jumpbox_private_key}"
echo "${BOSH_CA_CERT}" > "${BOSH_CA_CERT_PATH}"

chmod 0600 "${jumpbox_private_key}"
export BOSH_CA_CERT_PATH

ssh_cmd="ssh -i ${jumpbox_private_key} -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"

sshuttle -r "${JUMPBOX_USERNAME}@${JUMPBOX_HOST}" "10.0.0.0/16" -e "${ssh_cmd}" &
sshuttle_pid="$!"

trap "kill ${sshuttle_pid}" EXIT

pushd git-cfcr-etcd-release/src/bbr
  source ../../.envrc
  dep ensure
  ginkgo
popd

