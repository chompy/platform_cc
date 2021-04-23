#!/usr/bin/env bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
INSTALL_PATH=~/.pcc

# copy scripts
cp $SCRIPTPATH/uninstall.sh $INSTALL_PATH/pcc_uninstall

# build
bash $SCRIPTPATH/build_local.sh