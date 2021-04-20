#!/usr/bin/env bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
MODULE_PATHS=(".." "../api/container" "../api/def" "../api/output" "../api/platformsh" "../api/project" "../api/router" "../api/config" "../cli")

for i in "${MODULE_PATHS[@]}"
do
    cd $SCRIPTPATH/$i
    go mod tidy
done