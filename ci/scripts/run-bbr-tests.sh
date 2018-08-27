#!/bin/bash

set -euo pipefail

tar xvf bbr-release/bbr-*.tar
mv releases/bbr /usr/local/bin/bbr

version="v6.8.4"
wget https://github.com/cloudfoundry/bosh-bootloader/releases/download/${version}/bbl-${version}_linux_x86-64 -O /usr/local/bin/bbl
chmod +x /usr/local/bin/bbl

pushd bbl-state
  eval "$(bbl print-env)"
  jumpbox_host="$(bosh int --path /jumpbox_url <(bbl outputs))"
popd

BOSH_CA_CERT_PATH="$(mktemp)"
echo "${BOSH_CA_CERT}" > "${BOSH_CA_CERT_PATH}"
export BOSH_CA_CERT_PATH

ETCD_CLIENT_CERT="$(mktemp)"
ETCD_KEY_FILE="$(mktemp)"
ETCD_CA="$(mktemp)"

credhub get -n /bosh-etcd-bosh/bbr-etcd-single-node/tls_etcdctl --output-json | jq -r .value.certificate > "$ETCD_CLIENT_CERT"
credhub get -n /bosh-etcd-bosh/bbr-etcd-single-node/tls_etcdctl --output-json | jq -r .value.private_key > "$ETCD_KEY_FILE"
credhub get -n /bosh-etcd-bosh/bbr-etcd-single-node/tls_etcd --output-json | jq -r .value.ca > "$ETCD_CA"

ETCD_ENDPOINT="https://$(bosh vms --json | jq -r .Tables[0].Rows[0].ips):2379"
export ETCD_ENDPOINT ETCD_CLIENT_CERT ETCD_KEY_FILE ETCD_CA

ssh_cmd="ssh -i ${JUMPBOX_PRIVATE_KEY} -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"
sshuttle -r "jumpbox@${jumpbox_host}" "10.0.0.0/16" -e "${ssh_cmd}" &
sshuttle_pid="$!"

trap "kill ${sshuttle_pid}" EXIT

pushd git-cfcr-etcd-release
  source .envrc
popd

pushd git-cfcr-etcd-release/src/bbr
  dep ensure
  ginkgo
popd

