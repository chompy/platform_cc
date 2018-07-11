from .mariadb import MariaDbService
from .memcached import MemcachedService
from .rabbitmq import RabbitMqService
from .athenapdf import AthenaPdfService

""" Map service names to their service class. """
SERVICES_MAP = {
    "mariadb"       : MariaDbService,
    "mysql"         : MariaDbService,
    "memcached"     : MemcachedService,
    "rabbitmq"      : RabbitMqService,
    "athenapdf"     : AthenaPdfService
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