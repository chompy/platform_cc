#!/usr/bin/env bash

INSTALL_PATH=~/.pcc
OLD_BIN_PATHS=(~/.local/bin /usr/local/bin)
OLD_PCC_FILES=("pcc" "pcc_send_log" "pcc_psh_sync" "pcc_update" "pcc_uninstall")
BASH_INIT_PATHS=(~/.bashrc ~/.bash_profile ~/.zshrc)
SED_ARGS="-i"
if [[ "$OSTYPE" == "darwin"* ]]; then
    SED_ARGS="-i .bak"
fi

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
delete_old_file() {
    for s in ${OLD_BIN_PATHS[@]}; do
        if [ -L $s/$1 ] || [ -f $s/$1 ]; then
            printf "> Remove '$1'..."
            set -e
            rm -f $s/$1
            if [[ ! "$?" = "0" ]]; then
                progress_failed 
            else          
                progress_success
            fi
        fi
    done
}

# warn user
echo "UNINSTALLING PLATFORM.CC IN 5 SECONDS..."
sleep 5

# remove main installation directory
if [ -d $INSTALL_PATH ]; then
    printf "> Remove main directory..."
    set -e
    rm -r $INSTALL_PATH
    if [[ ! "$?" = "0" ]]; then
        progress_failed
        echo "ERROR: Failed to remove main directory."      
        exit 1
    fi
    progress_success
fi

# delete old files
for f in ${OLD_PCC_FILES[@]}; do    
    delete_old_file $f   
done

# remove bash complete
printf "> Remove PATH and bash completion..."
for b in ${BASH_INIT_PATHS[@]}; do
    if [ -f $b ]; then
        sed $SED_ARGS "s/source.*\.pcc.*//g" $b
    fi
done
progress_success