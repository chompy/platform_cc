#!/bin/sh

# set php ini config
echo "date.timezone = UTC" > /usr/local/etc/php/conf.d/01-main.ini
echo "memory_limit = 512M" >> /usr/local/etc/php/conf.d/01-main.ini
ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php5-fpm
ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.0
ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.1-zts

# install composer
php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
php composer-setup.php --install-dir=/usr/local/bin
rm composer-setup.php
ln -s -f /usr/local/bin/composer.phar /usr/local/bin/composer
