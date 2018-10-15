#!/bin/sh

docker build -t registry.gitlab.com/contextualcode/platform_cc/php54-fpm ./php54-fpm
docker build -t registry.gitlab.com/contextualcode/platform_cc/php56-fpm ./php56-fpm
docker build -t registry.gitlab.com/contextualcode/platform_cc/php70-fpm ./php70-fpm
docker build -t registry.gitlab.com/contextualcode/platform_cc/php71-fpm ./php71-fpm
docker build -t registry.gitlab.com/contextualcode/platform_cc/php72-fpm ./php72-fpm
docker build -t registry.gitlab.com/contextualcode/platform_cc/golang-1-11 ./golang-1-11
