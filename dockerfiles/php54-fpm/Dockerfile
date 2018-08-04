FROM php:5.4-fpm
LABEL com.contextualcode.platformcc=""

# install dependencies via apt-get
RUN apt-get update && \
    apt-get install -y rsync git unzip cron python-pip python-dev \
        gem libyaml-0-2 libyaml-dev ruby ruby-dev less \
        libmcrypt4 libmcrypt-dev libicu52 libicu-dev libxslt1.1 libxslt1-dev \
        libfreetype6 libfreetype6-dev libjpeg62-turbo libjpeg62-turbo-dev \
        libpng12-0 libpng-dev libpcre3 libpcre3-dev libedit-dev && \
    apt-get clean

# install nodejs
RUN curl https://nodejs.org/dist/v8.11.2/node-v8.11.2-linux-x64.tar.xz -o node.tar.xz && \
    tar xf node.tar.xz && \
    rm node.tar.xz && \
    mv node-* /opt/nodejs && \
    ln -s -f /opt/nodejs/bin/* /usr/bin/ && \
    ln -s -f /usr/bin/node /usr/bin/nodejs

# update php configurations
RUN echo "date.timezone = UTC" > /usr/local/etc/php/conf.d/01-main.ini && \
    echo "memory_limit = 512M" >> /usr/local/etc/php/conf.d/01-main.ini && \
    ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php5-fpm && \
    ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.0 && \
    ln -s -f /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.1-zts

# install composer
RUN php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');" && \
    php composer-setup.php --install-dir=/usr/local/bin && \
    rm composer-setup.php && \
    ln -s -f /usr/local/bin/composer.phar /usr/local/bin/composer

# install default php extensions
RUN docker-php-ext-install bcmath intl xsl mysql mysqli pdo_mysql sockets exif mcrypt && \
    docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ && \
    docker-php-ext-install gd

# install nginx + conf
RUN curl -L https://nginx.org/download/nginx-1.14.0.tar.gz -o nginx.tar.gz && \
    tar xvfz nginx.tar.gz && \
    rm nginx.tar.gz && \
    cd nginx* && \
    git clone --recursive https://github.com/google/ngx_brotli.git && \
    ./configure --with-http_realip_module --with-http_gunzip_module --with-http_gzip_static_module --add-module=ngx_brotli && \
    make && \
    make install && \
    cd .. && \
    rm -rf nginx* 
COPY nginx.conf /usr/local/nginx/conf/nginx.conf
COPY fastcgi_params /usr/local/nginx/conf/fastcgi_params
