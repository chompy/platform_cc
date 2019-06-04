#!/bin/sh

docker build -t registry.gitlab.com/contextualcode/platform_cc/php54-fpm ./php -f ./php/Dockerfile-php54
docker build -t registry.gitlab.com/contextualcode/platform_cc/php56-fpm ./php -f ./php/Dockerfile-php56
docker build -t registry.gitlab.com/contextualcode/platform_cc/php70-fpm ./php -f ./php/Dockerfile-php70
docker build -t registry.gitlab.com/contextualcode/platform_cc/php71-fpm ./php -f ./php/Dockerfile-php71
docker build -t registry.gitlab.com/contextualcode/platform_cc/php72-fpm ./php -f ./php/Dockerfile-php72
docker build -t registry.gitlab.com/contextualcode/platform_cc/php73-fpm ./php -f ./php/Dockerfile-php73
docker build -t registry.gitlab.com/contextualcode/platform_cc/golang-1-11 ./golang -f ./golang/Dockerfile-golang111
docker build -t registry.gitlab.com/contextualcode/platform_cc/python-3-7 ./python -f ./python/Dockerfile-python37
