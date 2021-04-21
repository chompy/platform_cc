#!/usr/bin/env bash

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PCC_BIN_NAME="pcc"
DATA_PATH=$SCRIPTPATH/../pkg/deb
VERSION=`cat $SCRIPTPATH/../version`

# use dev version if -d flag set
while getopts 'd' flag; do
    if [ "$flag" = "d" ]; then
        VERSION="dev"
    fi
done

# create build path and define deb path
DEB_PATH="$SCRIPTPATH/../build/$VERSION/platform_cc_$VERSION.deb"
mkdir -p `dirname $DEB_PATH`

# update version in control
sed -i "s/Version: .*/Version: $VERSION/g" $SCRIPTPATH/../pkg/deb/DEBIAN/control

# build
cd $SCRIPTPATH/..
go build -ldflags "-X main.version=$VERSION" -o $DATA_PATH/usr/bin/$PCC_BIN_NAME
dpkg-deb --build $DATA_PATH $DEB_PATH