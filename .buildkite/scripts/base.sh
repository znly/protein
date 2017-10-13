#!/bin/bash

set -eux

if [ "$#" -ne 1 ]; then
    echo "Missing image name"
    exit 1
fi

DKR_IMG_PATH=$1
GH_KEY="${HOME}/.ssh/id_rsa"
CI_SA="${HOME}/.ci.json"

cp "${GH_KEY}" "${CI_SA}" .
docker build -t "${DKR_IMG_PATH}" --build-arg GS_PATH="${GS_PATH}" -f Dockerfile.ci .
docker push "${DKR_IMG_PATH}"
