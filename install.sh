#!/usr/bin/env bash


BASE_URL="https://platform-cc-releases.s3.amazonaws.com"

# fetch current version
echo "Fetch latest version..."
VERSION=`curl -s $BASE_URL/version`
if [ -z "$VERSION" ]; then
    echo "ERROR: Could not determine latest version."
    exit 1
fi

# determine os
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    BASE_URL="$BASE_URL/$VERSION/linux_amd64"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    BASE_URL="$BASE_URL/$VERSION/darwin_amd64"
else
    BASE_URL="$BASE_URL/$VERSION/windows_amd64"
fi

# create local install
echo "Download version $VERSION..."
mkdir -p ~/.local/bin
rm -f ~/.local/bin/pcc
curl -s -l -o ~/.local/bin/pcc $BASE_URL
chmod +x ~/.local/bin/pcc
echo "DONE."