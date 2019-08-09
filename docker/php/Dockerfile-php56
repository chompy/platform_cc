FROM php:5.6-fpm
LABEL com.contextualcode.platformcc=""
COPY docker/php/conf/*.sh /opt/
COPY docker/php/conf/php-fpm.conf /etc/php/7.0/fpm/php-fpm.conf
RUN PHP_VER=5 /opt/install_dependencies.sh && \
    PHP_VER=5 /opt/config_php.sh && \
    PHP_VER=5 /opt/build_ext.sh
COPY docker/php/conf/nginx.conf /usr/local/nginx/conf/nginx.conf
COPY docker/php/conf/fastcgi_params /usr/local/nginx/conf/fastcgi_params