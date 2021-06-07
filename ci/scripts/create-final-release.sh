#!/usr/bin/env bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat cfcr-etcd-release-version/version)
git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"

cp -r git-cfcr-etcd-release/. git-cfcr-etcd-release-output

cd git-cfcr-etcd-release-output

cat <<EOF > "config/private.yml"
blobstore:
  provider: gcs
  options:
    credentials_source: static
    json_key: '${GCS_JSON_KEY}'
EOF

bosh create-release --final --version=${version} --sha2 --tarball ../cfcr-etcd-release/cfcr-etcd-release-${version}.tgz

echo "v${version}" >../cfcr-etcd-release/name
echo "" > ../cfcr-etcd-release/body

git checkout -b tmp/release
git add .
git commit -m "Final release for v${version}"
git tag -a "v${version}" -m "Tag for version v${version}"
git checkout "${BRANCH:-main}"
git merge tmp/release -m "Merge release branch for v${version}"
git branch -d tmp/release
