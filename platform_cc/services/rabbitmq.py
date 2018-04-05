from .base import BasePlatformService

class RabbitMqService(BasePlatformService):
    """
    Handler for Rabbitmq services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "rabbitmq:3.5"           : "rabbitmq:3",
        "rabbitmq:3.5"           : "rabbitmq:3"
    }

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerVolumes(self):
        volume = self.getVolume()
        return {
            volume.name : {
                "bind" : "/var/lib/rabbitmq",
                "mode" : "rw"
            }
        }

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"]["rabbitmq"] = {
            "host"          : self.getContainerName(),
            "ip"            : data.get("ip", ""),
            "scheme"        : "amqp",
            "port"          : 5672,
            "username"      : "guest",
            "password"      : "guest"
        }
        return data
