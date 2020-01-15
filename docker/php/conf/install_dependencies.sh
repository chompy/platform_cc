#!/bin/sh

# node.js
curl https://nodejs.org/dist/v10.16.0/node-v10.16.0-linux-x64.tar.gz -o node.tar.gz
tar xfz node.tar.gz
rm node.tar.gz
mv node-* /opt/nodejs
ln -s -f /opt/nodejs/bin/* /usr/bin/
ln -s -f /usr/bin/node /usr/bin/nodejs

# apt
mkdir -p /usr/share/man/man1
apt-get update
apt-get install -y rsync git unzip cron python-pip python-dev \
    gem libyaml-0-2 libyaml-dev ruby ruby-dev less nano libmemcached-dev  \
    libmcrypt4 libmcrypt-dev libxslt1.1 libxslt1-dev zlib1g-dev\
    libfreetype6 libfreetype6-dev libjpeg62-turbo libjpeg62-turbo-dev \
    libpcre3 libpcre3-dev libedit-dev gnupg apt-transport-https \
    imagemagick libmagickcore-dev libmagickwand-dev \
    libicu-dev libpng-dev \
    advancecomp jpegoptim libjpeg-turbo-progs optipng pngcrush
if [ "$PHP_VER" = "5" ]; then
    apt-get install -y default-jdk ant libcommons-lang3-java libbcprov-java
fi
if [ "$PHP_VER" = "7" ]; then
    apt-get install -y default-jdk-headless ant libcommons-lang3-java libbcprov-java
    ln -s /usr/lib/x86_64-linux-gnu/libicuuc.so /usr/lib/x86_64-linux-gnu/libicuuc.so.57
    ln -s /usr/lib/x86_64-linux-gnu/libicui18n.so /usr/lib/x86_64-linux-gnu/libicui18n.so.57
fi
apt-get clean

# ruby compass+sass
gem install compass sass
ln -s /usr/local/bin/compass /usr/bin/compass
ln -s /usr/local/bin/sass /usr/bin/sass

# yarn
cd /tmp
curl -L -o yarn.tgz https://yarnpkg.com/latest.tar.gz
tar xfz yarn.tgz
rm yarn.tgz
mkdir -p /opt/yarn
mv yarn-*/* /opt/yarn/
rm -rf yarn-*
ln -s /opt/yarn/bin/yarn /usr/bin/yarn

# nginx
curl -L https://nginx.org/download/nginx-1.17.1.tar.gz -o nginx.tar.gz
tar xfz nginx.tar.gz
rm nginx.tar.gz
cd nginx*
git clone --recursive https://github.com/google/ngx_brotli.git
./configure --with-http_realip_module --with-http_gunzip_module --with-http_gzip_static_module --add-module=ngx_brotli
make && \
make install && \
cd .. && \
rm -rf nginx* 
