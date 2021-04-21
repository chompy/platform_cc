#!/usr/bin/env bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PCC_BIN_NAME="pcc"
INSTALL_PATH=~/.pcc
if [ -f "/usr/bin/$PCC_BIN_NAME" ]; then
    INSTALL_PATH=/usr/bin
fi

# build
cd $SCRIPTPATH/..
go build -o $INSTALL_PATH/$PCC_BIN_NAME

# complete out, let user know where pcc was installed
PATH_TO=$INSTALL_PATH/$PCC_BIN_NAME
printf "\e[32mBUILT AT \e[0m$PATH_TO\e[0m\n"
