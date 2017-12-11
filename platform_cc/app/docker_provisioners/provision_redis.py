from __future__ import absolute_import
from .provision_base import DockerProvisionBase

class DockerProvision(DockerProvisionBase):

    """ Provision a redis container. """

    def getServiceRelationship(self):
        return [
            {
                "host" : self.container.attrs.get("NetworkSettings", {}).get("IPAddress", ""),
                "scheme" : "redis",
                "port" : "6379",
            }
        ]