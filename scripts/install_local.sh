#!/usr/bin/env bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
INSTALL_PATH=~/.local/bin
MAC_INSTALL_PATH=/usr/local/bin

# change install path for mac
if [[ "$OSTYPE" == "darwin"* ]]; then
    INSTALL_PATH="$MAC_INSTALL_PATH"
fi

# copy scripts
cp $SCRIPTPATH/platform_sh_clone.sh $INSTALL_PATH/pcc_psh_sync
cp $SCRIPTPATH/send_log.sh $INSTALL_PATH/pcc_send_log
cp $SCRIPTPATH/install.sh $INSTALL_PATH/pcc_update
cp $SCRIPTPATH/uninstall.sh $INSTALL_PATH/pcc_uninstall

# build
bash $SCRIPTPATH/build_local.sh