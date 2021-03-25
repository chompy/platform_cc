#!/usr/bin/env bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
INSTALL_PATH=~/.pcc

# copy scripts
cp $SCRIPTPATH/platform_sh_clone.sh $INSTALL_PATH/pcc_psh_sync
cp $SCRIPTPATH/send_log.sh $INSTALL_PATH/pcc_send_log
cp $SCRIPTPATH/install.sh $INSTALL_PATH/pcc_update
cp $SCRIPTPATH/uninstall.sh $INSTALL_PATH/pcc_uninstall

# build
bash $SCRIPTPATH/build_local.sh