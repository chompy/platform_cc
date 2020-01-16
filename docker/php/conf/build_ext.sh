#!/bin/sh

# install default extensions
docker-php-ext-install -j$(nproc) bcmath intl xsl mysqli pdo_mysql sockets exif zip
if [ "$PHP_VER" = "5" ]; then
    docker-php-ext-install -j$(nproc) mysql mcrypt
fi
docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/
docker-php-ext-install -j$(nproc) gd

# precompile extensions
build_ext()
{
    cd /tmp
    curl -o $2.tgz $1
    tar xvfz $2.tgz
    tar -tf $2.tgz | head -1 | cut -f1 -d"/"
    rm $2.tgz
    cd ${2}*
    phpize && ./configure && make && make install
    cd /tmp
    rm -rf $2*
}
if [ "$PHP_VER" = "5" ]; then
    build_ext "https://pecl.php.net/get/memcached-2.2.0.tgz" "memcached"
    build_ext "https://pecl.php.net/get/xdebug-2.5.5.tgz" "xdebug"
    build_ext "https://pecl.php.net/get/igbinary-2.0.8.tgz" "igbinary"
fi
if [ "$PHP_VER" = "7" ]; then
    build_ext "https://pecl.php.net/get/memcached-3.1.3.tgz" "memcached"
    build_ext "https://pecl.php.net/get/xdebug-2.7.2.tg" "xdebug"
    build_ext "https://pecl.php.net/get/igbinary-3.0.1.tgz" "igbinary"
fi
build_ext "https://pecl.php.net/get/redis-4.3.0.tgz" "redis"
build_ext "https://pecl.php.net/get/imagick-3.4.4.tgz" "imagick"

