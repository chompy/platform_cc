#!/bin/sh

# Build Platform.CC specific containers.

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

# router
docker build -t "contextualcode/platform_cc_router" -f "$SCRIPTPATH/../container/router/Dockerfile" "$SCRIPTPATH/../container/router"
docker push "contextualcode/platform_cc_router"