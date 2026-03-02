#!/usr/bin/env bash

set -e

k8s_version=v1.24.2
arch=amd64

if [[ "$OSTYPE" == "linux-gnu" ]]; then
  os="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  os="darwin"
else
  echo "OS '$OSTYPE' not supported." >&2
  exit 1
fi

root=$(cd "$(dirname "$0")"/..; pwd)
output_dir="$root"/_out
archive_name="envtest-$k8s_version-$os-$arch.tar.gz"
archive_file="$output_dir/$archive_name"
archive_url="https://github.com/kubernetes-sigs/controller-tools/releases/download/envtest-$k8s_version/$archive_name"

mkdir -p "$output_dir"
curl -sL "$archive_url" -o "$archive_file"
tar -zxf "$archive_file" -C "$output_dir/"
