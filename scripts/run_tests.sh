#!/usr/bin/env bash
set -e

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
TEST_PACKAGES=("api/container" "api/def" "api/project" "api/router" "api/output")

for p in ${TEST_PACKAGES[@]}; do
    echo "> $p..."
    echo ""
    cd $SCRIPTPATH/../$p
    go test -race -short -timeout 30s -coverprofile .pcc_coverage ./...
    rm .pcc_coverage
    echo ""
done
