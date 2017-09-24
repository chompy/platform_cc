import os
import difflib
import io
import hashlib
import docker
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a solr container. """

    def getEnvironmentVariables(self):
        return {
            "SOLR_HOME" : "/data",
            "INIT_SOLR_HOME" : "yes"
        }

    def getVolumes(self):
        config = self.appConfig.getConfiguration()
        volumeKey = "%s_%s_%s_data" % (
            DockerProvisionBase.DOCKER_VOLUME_NAME_PREFIX,
            self.appConfig.projectHash[:6],
            self.appConfig.getName()
        )
        volumes = {
            volumeKey : {
                "bind" : "/data",
                "mode" : "rw"
            }
        }
        try:
            self.dockerClient.volumes.get(volumeKey)
        except docker.errors.NotFound:
            self.dockerClient.volumes.create(volumeKey)
        container = self.dockerClient.containers.run(
            self.appConfig.getDockerImage(),
            name="%s_%s_solr_provisioner" % (
                DockerProvisionBase.DOCKER_VOLUME_NAME_PREFIX,
                self.appConfig.projectHash[:6]
            ),
            command="chown -R solr:solr /data",
            detach=True,
            volumes=volumes,
            hostname="%s_%s_solr_provisioner" % (
                DockerProvisionBase.DOCKER_VOLUME_NAME_PREFIX,
                self.appConfig.projectHash[:6]
            ),
            stdin_open=True,
            user="root"
        )
        container.wait()
        container.remove()
        return volumes

    def preBuild(self):
        config = self.appConfig.getConfiguration()
        cores = config.get("cores", {}).keys()
        if config.get("core_config", None):
            cores.append("collection1")
        cmds = []
        # create solr cores
        for core in cores:
            cmds.append({
                "cmd" : "/opt/solr/bin/solr create_core -c '%s'" % core,
                "desc" : "Create SOLR core '%s' if not exist." % core,
                "user" : "solr"
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