from .base import BasePlatformService
import hashlib
import base36
import docker
import time
import requests


class SolrService(BasePlatformService):
    """
    Handler for Solr service.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "solr:4.10":            "klabs/solr:ezfind",
        # "solr:6.3":            "solr:6.3-alpine",
        # "solr:6.6":            "solr:6.6-alpine"
    }

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerVolumes(self):
        return {
            self.getVolumeName(): {
                "bind": "/solr_data",
                "mode": "rw"
            }
        }

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        endpoints = self.config.get("endpoints", {})
        # for solr version 4 automatically add "collection1" endpoint
        # see https://docs.platform.sh/configuration/services/solr.html#solr-4
        versionNo = self.getType().split(":")[1][0]
        if versionNo == "4":
            endpoints = {
                "solr": {
                    "core": "collection1"
                }
            }
        for name, config in endpoints.items():
            data["platform_relationships"][name.strip()] = {
                "host":           self.getContainerName(),
                "ip":             data.get("ip", ""),
                "port":           8983,
                "path":           "solr/%s" % config.get("core", "default"),
                "scheme":         "solr"
            }
        return data

    def start(self):
        BasePlatformService.start(self)
        container = self.getContainer()
        if not container:
            return

        # copy configsets

        # copy config
        container.runCommand(
            """
            mkdir /solr_conf
            """
        )
        data = BasePlatformService.getServiceData(self)
        containerIp = data.get("ip", "")

        # configure cores
        # /solr/admin/cores?action=CREATE&name=%s&instanceDir=/opt/solr/data/%s&config=solrconfig.xml&dataDir=data
        # try:
        #    r = request.get("http://%s" % data.get("ip", ""))
        # except requests.exceptions.ConnectionError:
        #    pass

        # TODO
