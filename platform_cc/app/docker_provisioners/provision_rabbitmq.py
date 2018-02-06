from __future__ import absolute_import
from .provision_base import DockerProvisionBase

class DockerProvision(DockerProvisionBase):

    """ Provision a rabbitmq container. """

    def getServiceRelationship(self):
        return [
            {
                "host" : self.container.attrs.get("Config", {}).get("Hostname", ""),
                "ip" : self.container.attrs.get("NetworkSettings", {}).get("IPAddress", ""),
                "scheme" : "amqp",
                "port" : "5672",
                "username" : "guest",
                "password" : "guest"
            }
        ]