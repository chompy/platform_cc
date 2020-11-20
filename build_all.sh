#!/bin/bash

VERSION=`cat version`
PLATFORMS=("darwin/amd64" "windows/amd64" "linux/amd64")
RELEASE_BUCKET="platform-cc-releases"

# prepare build directory
mkdir -p build/$VERSION
rm -rf build/$VERSION/*

# itterate each platform and build
for platform in "${PLATFORMS[@]}"; do
    IFS="/" read platform arch <<< "$platform"
    echo "$platform ($arch)..."
    echo " > BUILD."
    GOOS=$platform GOARCH=$arch go build \
        -ldflags "-X main.version=$VERSION" \
        -o "build/$VERSION/${platform}_${arch}"

    # upload to s3 release bucket
    if [ ! -z "$AWS_ACCESS_KEY_ID" ] && [ ! -z "$AWS_SECRET_ACCESS_KEY" ]; then
        echo " > UPLOAD."
        date=`date +%Y%m%d`
        dateFormatted=`date -R`
        fileName="$VERSION/${platform}_${arch}"
        relativePath="/${RELEASE_BUCKET}/${fileName}"
        contentType="application/octet-stream"
        stringToSign="PUT\n\n${contentType}\n${dateFormatted}\n${relativePath}"
        s3AccessKey="$AWS_ACCESS_KEY_ID"
        s3SecretKey="$AWS_SECRET_ACCESS_KEY"
        signature=`echo -en ${stringToSign} | openssl sha1 -hmac ${s3SecretKey} -binary | base64`
        curl -X PUT -T "build/${fileName}" \
        -H "Host: ${RELEASE_BUCKET}.s3.amazonaws.com" \
        -H "Date: ${dateFormatted}" \
        -H "Content-Type: ${contentType}" \
        -H "Authorization: AWS ${s3AccessKey}:${signature}" \
        http://${RELEASE_BUCKET}.s3.amazonaws.com/${fileName}
        echo " > DONE."
    fi
done

# upload version
date=`date +%Y%m%d`
dateFormatted=`date -R`
fileName="version"
relativePath="/${RELEASE_BUCKET}/version"
contentType="application/octet-stream"
stringToSign="PUT\n\n${contentType}\n${dateFormatted}\n${relativePath}"
s3AccessKey="$AWS_ACCESS_KEY_ID"
s3SecretKey="$AWS_SECRET_ACCESS_KEY"
signature=`echo -en ${stringToSign} | openssl sha1 -hmac ${s3SecretKey} -binary | base64`
curl -X PUT -T "${fileName}" \
-H "Host: ${RELEASE_BUCKET}.s3.amazonaws.com" \
-H "Date: ${dateFormatted}" \
-H "Content-Type: ${contentType}" \
-H "Authorization: AWS ${s3AccessKey}:${signature}" \
http://${RELEASE_BUCKET}.s3.amazonaws.com/${fileName}