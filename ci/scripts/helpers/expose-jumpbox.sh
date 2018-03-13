#!/bin/bash

function expose_jumpbox {
  jumpbox_key_path="${1}"
  jumpbox_ssh_key="${2}"
  jumpbox_url="${3}"
  echo "${jumpbox_ssh_key}" > "${jumpbox_key_path}/jumpbox.key"
  chmod 600 "${jumpbox_key_path}/jumpbox.key"
  export BOSH_ALL_PROXY="ssh+socks5://${jumpbox_url}?private-key=${jumpbox_key_path}/jumpbox.key"
}
