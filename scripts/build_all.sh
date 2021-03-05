#!/usr/bin/env bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
ROOTPATH="$SCRIPTPATH/.."
VERSION=`cat $ROOTPATH/version`
PLATFORMS=("darwin/amd64" "windows/amd64" "linux/amd64")
RELEASE_BUCKET="platform-cc-releases"

# prepare build directory
mkdir -p $ROOTPATH/build/$VERSION
rm -rf $ROOTPATH/build/$VERSION/*

# use dev version if -d flag set
while getopts 'd' flag; do
    if [ "$flag" = "d" ]; then
        VERSION="dev"
    fi
done

# s3 upload function
upload() {
    echo " > UPLOAD $1 >> $2"
    date=`date +%Y%m%d`
    dateFormatted=`date -R`
    fileName="$2"
    relativePath="/${RELEASE_BUCKET}/${fileName}"
    contentType="application/octet-stream"
    stringToSign="PUT\n\n${contentType}\n${dateFormatted}\n${relativePath}"
    s3AccessKey="$AWS_ACCESS_KEY_ID"
    s3SecretKey="$AWS_SECRET_ACCESS_KEY"
    signature=`echo -en ${stringToSign} | openssl sha1 -hmac ${s3SecretKey} -binary | base64`
    curl -X PUT -T "$1" \
    -H "Host: ${RELEASE_BUCKET}.s3.amazonaws.com" \
    -H "Date: ${dateFormatted}" \
    -H "Content-Type: ${contentType}" \
    -H "Authorization: AWS ${s3AccessKey}:${signature}" \
    http://${RELEASE_BUCKET}.s3.amazonaws.com/${fileName}
}

# itterate each platform and build
for platform in "${PLATFORMS[@]}"; do
    IFS="/" read platform arch <<< "$platform"
    echo "- $platform ($arch)..."
    echo " > BUILD."
    GOOS=$platform GOARCH=$arch go build \
        -ldflags "-X main.version=$VERSION" \
        -o "$ROOTPATH/build/$VERSION/${platform}_${arch}" \
        $ROOTPATH/
    # upload to s3 release bucket
    if [ ! -z "$AWS_ACCESS_KEY_ID" ] && [ ! -z "$AWS_SECRET_ACCESS_KEY" ]; then
        upload "$ROOTPATH/build/$VERSION/${platform}_${arch}" "$VERSION/${platform}_${arch}"
    fi
done
# upload files
if [ ! -z "$AWS_ACCESS_KEY_ID" ] && [ ! -z "$AWS_SECRET_ACCESS_KEY" ]; then
    if [ "$VERSION" != "dev" ]; then
        echo "- Upload additional files..."
        upload "$ROOTPATH/version" "version"
        upload "$ROOTPATH/scripts/send_log.sh" "send_log.sh"
        upload "$ROOTPATH/scripts/install.sh" "install.sh"
        upload "$ROOTPATH/scripts/uninstall.sh" "uninstall.sh"
        upload "$ROOTPATH/scripts/platform_sh_clone.sh" "platform_sh_clone.sh"
        upload "$ROOTPATH/scripts/.pcc.bashrc" ".pcc.bashrc"
    fi
fi