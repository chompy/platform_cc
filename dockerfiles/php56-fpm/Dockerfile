FROM php:5.6-fpm

# install dependencies via apt-get
RUN apt-get update && \
    apt-get install -y rsync git unzip cron python-pip python-dev \
        gem libyaml-0-2 libyaml-dev ruby ruby-dev nginx less nano \
        libmcrypt4 libmcrypt-dev libicu57 libicu-dev libxslt1.1 libxslt1-dev \
        libfreetype6 libfreetype6-dev libjpeg62-turbo libjpeg62-turbo-dev \
        libpng16-16 libpng-dev && \
    apt-get clean

# install nodejs
RUN curl https://nodejs.org/dist/v8.11.2/node-v8.11.2-linux-x64.tar.xz -o node.tar.xz && \
    tar xf node.tar.xz && \
    rm node.tar.xz && \
    mv node-* /opt/nodejs && \
    ln -s -f /opt/nodejs/bin/* /usr/bin/ && \
    ln -s -f /usr/bin/node /usr/bin/nodejs

# update php configurations
RUN sed -i "s/user = .*/user = web/g" /usr/local/etc/php-fpm.d/www.conf && \
    sed -i "s/group = .*/group = web/g" /usr/local/etc/php-fpm.d/www.conf && \
    echo "date.timezone = UTC" > /usr/local/etc/php/conf.d/01-main.ini && \
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
RUN docker-php-ext-install -j$(nproc) bcmath intl xsl mysql mysqli pdo_mysql sockets exif mcrypt && \
    docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ && \
    docker-php-ext-install -j$(nproc) gd

# create web user
RUN useradd -d /app -m -p secret~ --uid 1000 web && \
    usermod -a -G staff web && \
    mkdir -p /var/lib/gems && \
    chown -R web:web /var/lib/gems && \
    chown -R root:staff /usr/bin && \
    chmod -R g+rw /usr/bin

# install nginx conf
COPY nginx.conf /etc/nginx/nginx.conf
