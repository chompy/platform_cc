from .base import BasePlatformService
import hashlib
import base36
import docker
import time
import requests
import io
import os

class SolrService(BasePlatformService):
    """
    Handler for Solr service.
    Just a placeholder/dummy service for now.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        #"solr:4.10":            "klabs/solr:ezfind",
        #"solr:6.3":            "solr:6.3-alpine",
        #"solr:6.6":            "solr:6.6-alpine",
        "solr:3.6":            "busybox:1",
        "solr:4.10":           "busybox:1",
        "solr:6.3":            "solr:6.3-alpine",
        "solr:6.6":            "solr:6.6-alpine",
        "solr:7.6":            "busybox:1"        
    }

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerVolumes(self):
        return {
            self.getVolumeName(): {
                "bind": "/mnt/data",
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

        # link data
        self.runCommand(
            """
            bash -c '[ ! -f /mnt/data/solr.xml ] && cp -rf /opt/solr/server/solr/* /mnt/data/'
            rm -rf /opt/solr/server/solr
            ln -s /mnt/data /opt/solr/server/solr
            chown -R solr:solr /mnt/data
            """,
            user = "root"
        )

        # copy configsets
        for core, config in self.config.get("cores", {}).items():
            self.logger.info(
                "Create/update core '%s'.",
                core
            )
            solrConfData = config.get("conf_dir", {})
            if solrConfData:
                for path, data in solrConfData.items():
                    fObj = io.BytesIO(data.encode())
                    fullPath = os.path.join("/tmp/solr_conf/%s" % core, path)
                    self.runCommand(
                        """
                        mkdir -p "%s"
                        chown -R solr:solr /tmp/solr_conf/%s
                        """ % (
                            os.path.dirname(fullPath),
                            core
                        ),
                        user = "root"
                    )
                    self.uploadFile(
                        fObj,
                        os.path.join("/tmp/solr_conf/%s" % core, path)
                    )
            # TODO we might need to manually copy conf updates at each start up
            self.shell(
                """
                solr create_core -c %s -d %s
                """ % (
                    core,
                    "/tmp/solr_conf/%s" % core
                ),
                user = "solr"
            )

        # restart for config to take effect
        self.getContainer().restart()

