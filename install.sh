#!/usr/bin/env bash

BASE_URL="https://s.chompy.me/platform_cc"
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    BASE_URL="$BASE_URL/linux_amd64/platform_cc"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    BASE_URL="$BASE_URL/mac_amd64/platform_cc"
else
    BASE_URL="$BASE_URL/win_amd64/platform_cc"
fi
mkdir -p ~/.local/bin
rm -f ~/.local/bin/pcc
curl -l -o ~/.local/bin/pcc $BASE_URL
chmod +x ~/.local/bin/pcc