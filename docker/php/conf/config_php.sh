#!/bin/sh

# set php ini config
echo "date.timezone = UTC" > /usr/local/etc/php/conf.d/01-main.ini
echo "memory_limit = 512M" >> /usr/local/etc/php/conf.d/01-main.ini
ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php5-fpm
ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.0
ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.1-zts
ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.2-zts
ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.3-zts
ln -s -f /app /usr/local/lib/php/extensions/*/

# install composer
php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
php composer-setup.php --install-dir=/usr/local/bin
rm composer-setup.php
ln -s -f /usr/local/bin/composer.phar /usr/local/bin/composer

# add 'web' user
useradd -l -d /app -m -p secret~ --uid 1000 web
usermod -a -G staff web

# tweaks to allow 'web' access to various paths
mkdir -p /var/lib/gems
chown -R web:web /var/lib/gems
chown -R root:staff /usr/bin
chmod -R g+rw /usr/bin

# create php-fpm config structure like PSH
mkdir -p /etc/php/7.0/fpm
mkdir -p /etc/php/7.1-zts
mkdir -p /etc/php/7.2-zts
mkdir -p /etc/php/7.3-zts
mkdir -p /etc/php5
ln -s /etc/php/7.0/fpm /etc/php/7.1-zts/fpm
ln -s /etc/php/7.0/fpm /etc/php/7.2-zts/fpm
ln -s /etc/php/7.0/fpm /etc/php/7.3-zts/fpm
ln -s /etc/php/7.0/fpm /etc/php5/fpm
rm /usr/local/etc/php-fpm.conf
ln -s /etc/php/7.0/fpm/php-fpm.conf /usr/local/etc/php-fpm.conf
chmod +r /etc/php/7.0/fpm/*
chmod +r /usr/local/etc/*
