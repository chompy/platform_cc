#!/bin/sh

# node.js
curl https://nodejs.org/dist/v12.3.1/node-v12.3.1.tar.gz -o node.tar.gz
tar xvfz node.tar.gz
rm node.tar.gz
mv node-* /opt/nodejs
ln -s -f /opt/nodejs/bin/* /usr/bin/
ln -s -f /usr/bin/node /usr/bin/nodejs

# apt
apt-get update
apt-get install -y rsync git unzip cron python-pip python-dev \
    gem libyaml-0-2 libyaml-dev ruby ruby-dev less nano libmemcached-dev zlib1g-dev \
    libmcrypt4 libmcrypt-dev libicu57 libicu-dev libxslt1.1 libxslt1-dev \
    libfreetype6 libfreetype6-dev libjpeg62-turbo libjpeg62-turbo-dev \
    libpng16-16 libpng-dev libpcre3 libpcre3-dev libedit-dev gnupg apt-transport-https \
    imagemagick libmagickcore-dev libmagickwand-dev
apt-get update
if [ "$PHP_VER" = "7" ]; then
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.libxslt1
    apt-get update
    apt-get install -y yarn
fi
apt-get clean