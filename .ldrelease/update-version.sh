#!/bin/bash

set -ue

sed -i "/^version:/c version: ${LD_RELEASE_VERSION}" ./Chart.yaml
