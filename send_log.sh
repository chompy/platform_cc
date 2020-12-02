#!/bin/sh

PASTEBIN_UPLOAD="https://pastebin.chompy.me/documents"
PASTEBIN_DOWNLOAD="https://pastebin.chompy.me"
LOG_PATH=".platform_cc.log"

# ----
if [ ! -f $LOG_PATH ]; then
    echo "ERROR: Log not found."
    exit 1
fi

LOG_DATA=`cat $LOG_PATH`
RES=$(curl -s --data-binary "@$LOG_PATH" --header "Content-Type: text/plain" "$PASTEBIN_UPLOAD")
KEY=$(echo "$RES" | python3 -c 'import json,sys;obj=json.load(sys.stdin);print(obj["key"])')

echo ""
echo ">> $PASTEBIN_DOWNLOAD/$KEY"
echo ""