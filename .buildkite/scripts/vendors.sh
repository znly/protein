#!/bin/bash
set -eux

## GS_PATH is exported in buildkite.
test ${GS_PATH}

GLIDE_SHA=$(shasum ./Gopkg.lock | awk '{print $1}')
VENDOR_FILE="protein-${GLIDE_SHA}.tgz"

if gsutil stat "${GS_PATH}"/"${VENDOR_FILE}"; then
    gsutil cp "${GS_PATH}"/"${VENDOR_FILE}" .
    tar xf "${VENDOR_FILE}"
    rm "${VENDOR_FILE}"
else
    make -f Makefile deps
    tar czf "${VENDOR_FILE}" vendor
    gsutil -h "Cache-Control:private, max-age=0, no-transform" cp "${VENDOR_FILE}" "${GS_PATH}"/"${VENDOR_FILE}"
    rm "${VENDOR_FILE}"
fi
