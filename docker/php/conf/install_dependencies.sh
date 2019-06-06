#!/bin/sh

# node.js
curl https://nodejs.org/dist/v10.16.0/node-v10.16.0-linux-x64.tar.xz -o node.tar.xz
tar xf node.tar.xz
rm node.tar.xz
mv node-* /opt/nodejs
ln -s -f /opt/nodejs/bin/* /usr/bin/
ln -s -f /usr/bin/node /usr/bin/nodejs

# apt
apt-get update
apt-get install -y rsync git unzip cron python-pip python-dev \
    gem libyaml-0-2 libyaml-dev ruby ruby-dev less nano libmemcached-dev  \
    libmcrypt4 libmcrypt-dev libxslt1.1 libxslt1-dev zlib1g-dev\
    libfreetype6 libfreetype6-dev libjpeg62-turbo libjpeg62-turbo-dev \
     libpcre3 libpcre3-dev libedit-dev gnupg apt-transport-https \
    imagemagick libmagickcore-dev libmagickwand-dev
if [ "$PHP_VER" = "5" ]; then
    apt-get install -y libicu52 libicu-dev libpng12-0 libpng-dev
elif [ "$PHP_VER" = "7" ]; then
    apt-get install -y libicu57 libicu-dev libpng16-16 libpng-dev
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.libxslt1
    apt-get update
    apt-get install -y yarn
fi
apt-get clean

# nginx
curl -L https://nginx.org/download/nginx-1.14.0.tar.gz -o nginx.tar.gz
tar xvfz nginx.tar.gz
rm nginx.tar.gz
cd nginx*
git clone --recursive https://github.com/google/ngx_brotli.git
./configure --with-http_realip_module --with-http_gunzip_module --with-http_gzip_static_module --add-module=ngx_brotli
make && \
make install && \
cd .. && \
rm -rf nginx* 
