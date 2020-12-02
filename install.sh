#!/usr/bin/env bash

BASE_URL="https://platform-cc-releases.s3.amazonaws.com"
VERSION_URL="$BASE_URL/version"
SEND_LOG_URL="$BASE_URL/send_log.sh"
INSTALL_PATH=~/.local/bin
PCC_BIN_NAME="pcc"
SEND_LOG_BIN_NAME="pcc_send_log"

echo ""
printf "\e[33m================================\e[0m\n"
echo " PLATFORM.CC BY CONTEXTUAL CODE"
printf "\e[33m================================\e[0m\n"
echo ""

progress_success() {
    printf "\e[32mDONE\e[0m\n"
}
progress_error() {
    printf "\e[31mERROR\e[0m\n"
    printf "\n\e[31mERROR:\e[0m\n$1\n\n"
    exit 1
}

# fetch current version
printf "> Fetch current version number..."
VERSION=`curl -s --fail $VERSION_URL`
if [ -z "$VERSION" ]; then
    progress_error "Could not determine latest version."
fi
progress_success
printf "> Version \e[36m$VERSION\e[0m found.\n"

# fetch send log script
printf "> Fetch send log script..."
mkdir -p $INSTALL_PATH
curl -s --fail -o $INSTALL_PATH/$SEND_LOG_BIN_NAME "$SEND_LOG_URL"
if [ "$?" != "0" ]; then
    progress_error "Could not download send log script."
fi
chmod +x $INSTALL_PATH/$SEND_LOG_BIN_NAME
progress_success

# determine os
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    BASE_URL="$BASE_URL/$VERSION/linux_amd64"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    BASE_URL="$BASE_URL/$VERSION/darwin_amd64"
else
    BASE_URL="$BASE_URL/$VERSION/windows_amd64"
fi

# create local install
printf "> Download version \e[36m$VERSION\e[0m..."
rm -f $INSTALL_PATH/$PCC_BIN_NAME
curl -s --fail -o $INSTALL_PATH/$PCC_BIN_NAME $BASE_URL
if [ "$?" != "0" ]; then
    progress_error "Could not download version \e[36m$VERSION\e[0m."
fi
chmod +x $INSTALL_PATH/$PCC_BIN_NAME
progress_success
echo ""

# complete out, let user know where pcc was installed and provide readme
PATH_TO=`realpath $INSTALL_PATH/$PCC_BIN_NAME`
printf "\e[32mINSTALLED AT \e[0m$PATH_TO\e[0m\n"
printf "\e[32mSEE README AT \e[0mhttps://gitlab.com/contextualcode/platform_cc/-/blob/v2.0.x/README.md\e[0m\n\n"