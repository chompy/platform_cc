#!/usr/bin/env bash

realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}

BASE_URL="https://platform.cc/releases"
VERSION_URL="$BASE_URL/version"
DL_SCRIPTS=(
    "$BASE_URL/send_log.sh|pcc_send_log|send-log"
    "$BASE_URL/platform_sh_clone.sh|pcc_psh_sync|platform.sh-clone"
    "$BASE_URL/uninstall.sh|pcc_uninstall|uninstall"
    "$BASE_URL/install.sh|pcc_update|update"
)
INSTALL_PATH=~/.pcc
PCC_BIN_NAME="pcc"
BASH_INIT_PATHS=(~/.bashrc ~/.bash_profile)
SED_ARGS="-i"
if [[ "$OSTYPE" == "darwin"* ]]; then
    SED_ARGS="-i .bak"
fi

# set dev version if -d flag set
VERSION=""
while getopts 'dh' flag; do
    if [ "$flag" = "d" ]; then
        VERSION="dev"
    elif [ "$flag" = "h" ]; then
        echo "Install or update Platform.CC."
        echo ""
        echo "Flags:"
        echo "-h        Display help."
        echo "-d        Install lastest development release."
        exit 0
    fi
done

# display intro/title
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
if [ -z "$VERSION" ]; then
    printf "> Fetch current version number..."
    VERSION=`curl -s -L --fail $VERSION_URL`
    if [ -z "$VERSION" ]; then
        progress_error "Could not determine latest version."
    fi
    progress_success
    printf "> Version \e[36m$VERSION\e[0m found.\n"
fi

# make install dir
mkdir -p $INSTALL_PATH

# itterate and install scripts
for s in ${DL_SCRIPTS[@]}; do
    IFS='|' read -ra DL_SCRIPT <<< "$s"
    printf "> Fetch ${DL_SCRIPT[2]} script..."
    mkdir -p $INSTALL_PATH
    curl -s -L --fail -o $INSTALL_PATH/${DL_SCRIPT[1]} "${DL_SCRIPT[0]}"
    if [ "$?" != "0" ]; then
        progress_error "Could not download ${DL_SCRIPT[2]} script."
    fi
    chmod +x "$INSTALL_PATH/${DL_SCRIPT[1]}"
    progress_success
done

# bash completion
printf "> Add PATH and bash completion..."
curl -s -L --fail -o $INSTALL_PATH/.pcc.bashrc "$BASE_URL/.pcc.bashrc"
if [ "$?" != "0" ]; then
    progress_error "Could not download bashrc script."
fi
for b in ${BASH_INIT_PATHS[@]}; do
    if [ -f $b ]; then
        sed $SED_ARGS "s/source.*\.pcc.*//g" $b
        printf "source $INSTALL_PATH/.pcc.bashrc" >> $b
        source $INSTALL_PATH/.pcc.bashrc
    fi
done
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
curl -s -L --fail -o $INSTALL_PATH/$PCC_BIN_NAME $BASE_URL
if [ "$?" != "0" ]; then
    progress_error "Could not download version \e[36m$VERSION\e[0m."
fi
chmod +x $INSTALL_PATH/$PCC_BIN_NAME
progress_success
echo ""

# complete out, let user know where pcc was installed and provide readme
PATH_TO=`realpath $INSTALL_PATH/$PCC_BIN_NAME`
printf "\e[32mINSTALLED AT \e[0m$PATH_TO\e[0m\n"
printf "\e[32mSEE README AT \e[0mhttps://gitlab.com/contextualcode/platform_cc/-/blob/main/README.md\e[0m\n\n"
