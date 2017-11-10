import os
import base64
import json
import time
import docker

class PlatformVars:

    """ Handles storage and retrival of user defined variables. """

    VARS_DOCKER_IMAGE = "pcc_key_value_store"
    DOCKERFILE_PATH = "key_value_store"
    STORAGE_PATH = "/data"
    CONTAINER_CMD = "python /key_value_store.py"
    DOCKER_VOLUME_PREFIX = "pcc_"
    DOCKER_VOLUME_SUFFIX = "_vars_data"

    def __init__(self, projectHash):
        self.projectHash = projectHash
        self.dockerClient = docker.from_env()

    def _sanitizeKey(self, key):
        return key.replace(" ", "_").strip().lower()

    def _runCmd(self, action = "list", key = "", value =""):
        """ Send var command to key/value container. """
        try:
            self.dockerClient.images.get(
                self.VARS_DOCKER_IMAGE
            )
        except docker.errors.ImageNotFound:
            self.dockerClient.images.build(
                path=os.path.join(
                    os.path.dirname(__file__),
                    "../containers/%s" % (self.DOCKERFILE_PATH)
                ),
                tag=self.VARS_DOCKER_IMAGE
            )
        dockerVolumeKey = self.getVolumeKey()
        try:
            self.dockerClient.volumes.get(dockerVolumeKey)
        except docker.errors.NotFound:
            self.dockerClient.volumes.create(dockerVolumeKey)
        containerId = "pcc_%s_vars_%s" % (self.projectHash[:6], time.time())
        results = self.dockerClient.containers.run(
            self.VARS_DOCKER_IMAGE,
            name=containerId,
            command="%s %s%s%s" % (
                self.CONTAINER_CMD,
                action.strip().lower(),
                (" -k '%s'" % (self._sanitizeKey(key))) if key else "",
                (
                    " -v '%s'" % (
                        base64.b64encode(
                            bytes(
                                str(value).encode("utf-8")
                            )
                        ).decode("utf-8")
                    )
                ) if value else ""
            ),
            volumes={
                dockerVolumeKey : {
                    "bind" :    "/data",
                    "mode" :    "rw"
                }
            },
            remove=True
        )
        self.dockerClient.containers.prune({
            "label" : containerId
        })
        return results

    def getVolumeKey(self):
        return "%s%s%s" % (
            self.DOCKER_VOLUME_PREFIX,
            self.projectHash[:6],
            self.DOCKER_VOLUME_SUFFIX
        )

    def set(self, key, value):
        return self._runCmd(
            "set",
            key,
            value
        )

    def get(self, key):
        results = self._runCmd(
            "get",
            key
        )
        if results:
            results = base64.b64decode(results).decode("utf-8")
        return results

    def delete(self, key):
        return self._runCmd(
            "delete",
            key
        )

    def all(self):
        results = self._runCmd(
            "list"
        )
        if results:
            results = json.loads(
                results.decode("utf-8")
            )
            if results:
                for key in results:
                    results[key] = base64.b64decode(
                        results[key]
                    ).decode("utf-8")
        return results
