#!/bin/sh

docker build -t registry.gitlab.com/contextualcode/platform_cc/php56-fpm . -f ./docker/php/Dockerfile-php56
docker build -t registry.gitlab.com/contextualcode/platform_cc/php70-fpm . -f ./docker/php/Dockerfile-php70
docker build -t registry.gitlab.com/contextualcode/platform_cc/php71-fpm . -f ./docker/php/Dockerfile-php71
docker build -t registry.gitlab.com/contextualcode/platform_cc/php72-fpm . -f ./docker/php/Dockerfile-php72
docker build -t registry.gitlab.com/contextualcode/platform_cc/php73-fpm . -f ./docker/php/Dockerfile-php73
docker build -t registry.gitlab.com/contextualcode/platform_cc/php74-fpm . -f ./docker/php/Dockerfile-php74
docker build -t registry.gitlab.com/contextualcode/platform_cc/golang-1-11 . -f ./docker/golang/Dockerfile-golang111
docker build -t registry.gitlab.com/contextualcode/platform_cc/python-3-7 . -f ./docker/python/Dockerfile-python37
