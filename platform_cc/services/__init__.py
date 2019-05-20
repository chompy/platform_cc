"""
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
"""

from .mariadb import MariaDbService
from .memcached import MemcachedService
from .rabbitmq import RabbitMqService
from .athenapdf import AthenaPdfService
from .minio import MinioService
from .redis import RedisService
from .docker import DockerService
from .solr import SolrService
from .varnish import VarnishService

""" Map service names to their service class. """
SERVICES_MAP = {
    "mariadb"         : MariaDbService,
    "mysql"           : MariaDbService,
    "memcached"       : MemcachedService,
    "rabbitmq"        : RabbitMqService,
    "athenapdf"       : AthenaPdfService,
    "minio"           : MinioService,
    "redis"           : RedisService,
    "redis-persistent": RedisService,
    "docker"          : DockerService,
    "solr"            : SolrService,
    "varnish"         : VarnishService
}

def getService(project, config):
    """
    Get service handler from configuration.
    
    :param project: Dictionary containing project data
    :param config: Dictionary containing service configuration
    :return: Service handler object
    :rtype: .base.BasePlatformService
    """

    # validate config
    if not isinstance(config, dict):
        raise ValueError("Config parameter must be a dictionary (dict) object.")
    if "_type" not in config or "_name" not in config:
        raise ValueError("Config parameter is missing parameters required for a service.")
    serviceType = config["_type"].split(":")[0]
    if serviceType not in SERVICES_MAP:
        raise NotImplementedError(
            "No service handler available for '%s.'" % (
                config["_type"]
            )
        )

    # init service
    return SERVICES_MAP[serviceType](
        project,
        config
    )