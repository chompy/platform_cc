#!/usr/bin/env bash

INSTALL_PATH=~/.local/bin
MAC_INSTALL_PATH=/usr/local/bin
PCC_FILES=("pcc" "pcc_send_log" "pcc_psh_sync" "pcc_update", "pcc_uninstall")

progress_success() {
    printf "\e[32mDONE\e[0m\n"
}
progress_not_found() {
    printf "\e[35mMISSING\e[0m\n"
}
progress_failed() {
    printf "\e[31mFAILED\e[0m\n"
}
progress_error() {
    printf "\e[31mERROR\e[0m\n"
    printf "\n\e[31mERROR:\e[0m\n$1\n\n"
    exit 1
}
delete_file() {
    printf "> Remove '$1'..."
    if [[ ! -f $INSTALL_PATH/$1 ]]; then
        progress_not_found
    else
        set -e
        rm -f $INSTALL_PATH/$1
        if [[ "$?" = "0" ]]; then
            progress_success
        else
            progress_failed
        fi
    fi
}

# warn user
echo "UNINSTALLING PLATFORM.CC IN 5 SECONDS..."
sleep 5

# change install path for mac
if [[ "$OSTYPE" == "darwin"* ]]; then
    INSTALL_PATH="$MAC_INSTALL_PATH"
fi

# do
for f in ${PCC_FILES[@]}; do
    delete_file $f
done

# remove bash complete
sed -i "s/source.*\.pcc.*//g" ~/.bashrc