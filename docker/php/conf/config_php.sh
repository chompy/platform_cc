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

# add 'web' user
useradd -l -d /app -m -p secret~ --uid 1000 web
usermod -a -G staff web

# tweaks to allow 'web' access to various paths
mkdir -p /var/lib/gems
chown -R web:web /var/lib/gems
chown -R root:staff /usr/bin
chmod -R g+rw /usr/bin

# config
if [ -d /usr/local/etc/php-fpm.d ]; then
    sed -i "s/user = .*/user = web/g" /usr/local/etc/php-fpm.d/www.conf
    sed -i "s/group = .*/group = web/g" /usr/local/etc/php-fpm.d/www.conf
    sed -i "s/;listen.backlog.*/listen.backlog = 511/g" /usr/local/etc/php-fpm.d/www.conf
    sed -i "s/;listen.owner.*/listen.owner = web/g" /usr/local/etc/php-fpm.d/www.conf
    sed -i "s/;listen.group.*/listen.group = web/g" /usr/local/etc/php-fpm.d/www.conf
    sed -i "s/;listen.mode.*/listen.mode = 0660/g" /usr/local/etc/php-fpm.d/www.conf
    sed -i "s/listen.*/listen = \/run\/app.sock/g" /usr/local/etc/php-fpm.d/zz-docker.conf
else
    sed -i "s/user = .*/user = web/g" /usr/local/etc/php-fpm.conf
    sed -i "s/group = .*/group = web/g" /usr/local/etc/php-fpm.conf
    sed -i "s/;listen.backlog.*/listen.backlog = 511/g" /usr/local/etc/php-fpm.conf
fi