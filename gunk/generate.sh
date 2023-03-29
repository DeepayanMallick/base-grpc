#!/usr/bin/env bash

set -euxo pipefail

source ./devenv.sh

GOBIN=$(pwd)/bin go install \
        github.com/gunk/gunk

$(pwd)/bin/gunk generate ./...
