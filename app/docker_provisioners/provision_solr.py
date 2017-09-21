import os
import difflib
import io
import hashlib
import docker
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a solr container. """

    def getVolumes(self):
        config = self.appConfig.getConfiguration()
        volumes = {}
        cores = config.get("cores", {}).keys()
        if config.get("core_config", None):
            cores.append("collection1")
        for core in cores:
            coreVolumeKey = "%s_%s_%s_solrcore_%s" % (
                DockerProvisionBase.DOCKER_VOLUME_NAME_PREFIX,
                self.appConfig.projectHash[:6],
                self.appConfig.getName(),
                core
            )
            try:
                self.dockerClient.volumes.get(coreVolumeKey)
            except docker.errors.NotFound:
                self.dockerClient.volumes.create(coreVolumeKey)
            volumes[coreVolumeKey] = {
                "bind" : "/opt/solr/server/solr/%s" % core,
                "mode" : "rw"
            }
        return volumes

    def preBuild(self):
        cmds = []
        config = self.appConfig.getConfiguration()
        cores = config.get("cores", {}).keys()
        if config.get("core_config", None):
            cores.append("collection1")
        # create solr cores
        for core in cores:
            cmds.append({
                "cmd" : "solr create_core -c \"%s\"" % core,
                "desc" : "Create SOLR core '%s.'" % core
            })
        # add config dir
        # TODO
        self.runCommands(cmds)

    def getServiceRelationship(self):
        return [
            {
                "path" : "solr",
                "host" : self.container.attrs.get("NetworkSettings", {}).get("IPAddress", ""),
                "scheme" : "solr",
                "port" : "8983",
            }
        ]

    """def getServiceRelationship(self):
        config = self.appConfig.getConfiguration()
        endpoints = config.get("endpoints", {})
        relationships = []
        for endpointName in endpoints:
            endpoint = endpoints[endpointName]
            relationships.append({
                "path" : "solr/"
                "host" : self.container.attrs.get("Config", {}).get("Hostname", ""),
                "ip" : self.container.attrs.get("NetworkSettings", {}).get("IPAddress", ""),
                "password" : self.getPassword(endpointName),
                "path" : endpoint.get("default_schema", ""),
                "port" : "3306",
                "query": {
                    "is_master" : True # uncertain what makes this true (maybe first endpoint is master?)
                },
                "scheme" : "mysql",
                "username" : endpointName
            })
        return relationships
    """