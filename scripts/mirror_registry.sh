#!/usr/bin/env bash
FROM="docker.registry.platform.sh"
TO="registry.gitlab.com/contextualcode/platform_cc"

TAGS=(
    "php-5.4" "php-5.5" "php-5.6" "php-7.0" "php-7.1" "php-7.2" "php-7.3" "php-7.4"
    "mysql-5.5" "mysql-10.0" "mysql-10.1" "mysql-10.2" "mysql-10.3" "mysql-10.4"
    "mariadb-5.5" "mariadb-10.0" "mariadb-10.1" "mariadb-10.2" "mariadb-10.3" "mariadb-10.4"
    "solr-3.6" "solr-4.10" "solr-6.3" "solr-6.6" "solr-7.6" "solr-7.7" "solr-8.0" "solr-8.4" "solr-8.6"
    "varnish-5.1" "varnish-5.2" "varnish-6.0" "varnish-6.3"
    "redis-2.8" "redis-3.0" "redis-3.2" "redis-4.0" "redis-5.0" "redis-6.0"
    "memcached-1.4" "memcached-1.5" "memcached-1.6"
)

for tag in "${TAGS[@]}"
do
    docker pull $FROM/$tag
    docker tag $FROM/$tag $TO/$tag
    docker push $TO/$tag
    docker rmi $TO/$tag
done