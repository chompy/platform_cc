from .base import BasePlatformService
import io
import os
import json
import time
import tarfile

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
        "solr:7.6":            "solr:7.6",
        "solr:7.7":            "solr:7.7",
    }

    SAVE_PATH = "/mnt/data"
    CONF_PATH = "/mnt/config"

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerCommand(self):
        if self.getBaseImage() == "busybox:1":
            return []
        return ["solr-foreground", "-Dsolr.disable.shardsWhitelist=true"]

    def getContainerVolumes(self):
        return {
            self.getVolumeName(): {
                "bind": self.SAVE_PATH,
                "mode": "rw"
            }
        }

    def getVersionNo(self):
        versionNo = self.getType().split(":")[1][0]
        return versionNo

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        endpoints = self.config.get("endpoints", {})
        # for solr version 4 automatically add "collection1" endpoint
        # see https://docs.platform.sh/configuration/services/solr.html#solr-4
        versionNo = self.getVersionNo()
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

    def getSolrConfigData(self):
        """ Get SOLR config data files in tar.gz format. """
        tarIO = io.BytesIO()
        with tarfile.open(fileobj=tarIO, mode="w:gz") as tf:
            for core, config in self.config.get("cores", {}).items():
                solrConfData = config.get("conf_dir", {})
                if not solrConfData: continue
                for path, data in solrConfData.items():
                    tarInfo = tarfile.TarInfo("%s/%s" % (core, path))
                    tarInfo.size = len(data)
                    tf.addfile(tarInfo, io.BytesIO(data.encode()))
        tarIO.seek(0)
        return tarIO

    def getStartCommand(self):
        if self.getBaseImage() == "busybox:1":
            return ""
        return """
            bash -c '[ ! -f %s/solr.xml ] && cp -rf /opt/solr/server/solr/* %s/'
            rm -rf /opt/solr/server/solr
            ln -s %s /opt/solr/server/solr
            chown -R solr:solr %s
        """ % (
            self.SAVE_PATH,
            self.SAVE_PATH,
            self.SAVE_PATH,
            self.SAVE_PATH
        )

    def getCreateCoresCommand(self):
        if self.getBaseImage() == "busybox:1":
            return ""
        """ Get shell commands needed to create SOLR cores. """
        output = ""
        for core in self.config.get("cores", {}):
            savePath = "%s/%s" % (self.SAVE_PATH, core)
            confPath = "%s/%s" % (self.CONF_PATH, core)
            output += """
                if [ -f %s/solrconfig.xml ]; then
                    if ! grep "schemaFactory class" %s/solrconfig.xml > /dev/null; then
                        sed -i 's/<codecFactory class/<schemaFactory class="ClassicIndexSchemaFactory"><\/schemaFactory><codecFactory class/' %s/solrconfig.xml
                    fi
                fi
                if [ ! -d %s ]; then
                    solr create_core -c %s -d %s
                fi
            """ % (
                confPath,
                confPath,
                confPath,
                savePath,
                core,
                confPath
            )
        return output

    def start(self):
        BasePlatformService.start(self)
        # start up command
        self.runCommand(
            self.getStartCommand(),
            user="root"
        )
        # upload config
        configData = self.getSolrConfigData()
        self.runCommand(
            "mkdir -p %s" % self.CONF_PATH,
            user="root"
        )
        self.uploadFile(
            configData,
            "%s/config.tar.gz" % self.CONF_PATH
        )
        self.runCommand(
            """
            cd %s
            tar xvfz %s/config.tar.gz
            chown -R solr:solr %s/*
            """ % (
                self.CONF_PATH,
                self.CONF_PATH,
                self.CONF_PATH
            )
        )
        # wait for solr to become ready
        if self.getBaseImage() == "busybox:1":
            return
        exitCode = 1
        while exitCode != 0:
            (exitCode, _) = self.getContainer().exec_run(
                "nc -z 127.0.0.1 8983"
            )
            time.sleep(.35)
        # create cores
        self.runCommand(
            self.getCreateCoresCommand(),
            user="solr"
        )
        # restart for config to take effect
        self.getContainer().restart()

