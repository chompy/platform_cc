"""
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
"""

import os
import io
import time
import tarfile
import docker
import logging
import json
import tempfile
from dockerpty import PseudoTerminal, ExecOperation
from ..exception.state_error import StateError
from ..exception.container_command_error import ContainerCommandError

class Container:
    """
    Base class for all Docker container instances.
    """

    """ Name of repository for all committed application images. """
    COMMIT_REPOSITORY_NAME = "platform_cc"

    """ Name prefix to use for all service containers. """
    CONTAINER_NAME_PREFIX = "pcc_"

    """ Prefix to use for all Docker labels. """
    LABEL_PREFIX = "com.contextualcode.platformcc"


    def __init__(self, project, name, dockerClient=None):
        """
        Constructor.

        :param project: Project data
        :param name: Name of this container
        :param dockerClient: Docker client
        """
        self.project = dict(project)
        self.name = str(name)
        self.docker = dockerClient
        if not self.docker:
            self.docker = docker.from_env(
                timeout=1800  # 30 minutes
            )
        self.logger = logging.getLogger(__name__)
        self._container = None
        self._hasCommitImage = None

    def getName(self):
        """
        Get name of container service/application/etc.

        :return: Service or application name
        :rtype: str
        """
        return self.name

    def getDockerImage(self):
        """
        Get docker image name for container.

        :return: Docker image name
        :rtype: str
        """
        commitImage = self.getCommitImage()
        if commitImage and self._hasCommitImage is None:
            try:
                self.docker.images.get(commitImage)
                self._hasCommitImage = True
            except docker.errors.ImageNotFound:
                self._hasCommitImage = False
        if commitImage and self._hasCommitImage is True:
            return commitImage
        return self.getBaseImage()

    def getBaseImage(self):
        """
        Get base Docker image name to use for application
        prior to build.

        :return: Docker image name
        :rtype: str
        """
        return "busybox:1"

    def getCommitImage(self):
        """
        Get name of committed Docker image to use
        instead of main image if it exists.

        :return: Committed Docker image name
        :rtype: str
        """
        return "%s:%s_%s" % (
            self.COMMIT_REPOSITORY_NAME,
            self.getName(),
            self.project.get("short_uid")
        )

    @staticmethod
    def staticGetContainerName(project, name):
        """
        Get name of docker container.
        :param project: Project data
        :param name: Container base name
        :return: Docker container name
        :rtype: str
        """
        return "%s%s_%s" % (
            Container.CONTAINER_NAME_PREFIX,
            project.get("short_uid"),
            name
        )

    def getContainerName(self):
        """
        Get name of docker container.

        :return: Docker container name
        :rtype: str
        """
        return Container.staticGetContainerName(
            self.project,
            self.name
        )

    @staticmethod
    def staticGetNetworkName(project):
        """
        Get name of network to use with Docker container.

        :param project: Project data
        :return: Network name
        :rtype: str
        """
        return "%s%s" % (
            Container.CONTAINER_NAME_PREFIX,
            project.get("short_uid")
        )

    def getNetworkName(self):
        """
        Get name of network to use with Docker container.

        :return: Network name
        :rtype: str
        """
        return Container.staticGetNetworkName(
            self.project
        )

    def getVolumeName(self, name=""):
        """
        Get name of volume for use with this Docker container.

        :return: Volume name
        :rtype: str
        """
        volumeId = "%s%s_%s" % (
            self.CONTAINER_NAME_PREFIX,
            self.project.get("short_uid"),
            self.getName()
        )
        if name:
            volumeId = "%s_%s" % (
                volumeId,
                str(name)
            )
        return volumeId

    def getContainerCommand(self):
        """
        Get command to run to start service inside container.

        :return: Command or None if default command should be used
        """
        return None

    def getContainerEnvironmentVariables(self):
        """
        Get dictionary of environment variables to set inside container.

        :return: Dictionary of environment variables.
        :rtype: dict
        """
        return {}

    def getContainerHosts(self):
        """
        Addtional hostnames to resolve inside the container,
        as a mapping of hostname to IP address.

        :return: Dictionary of hostname to ip map
        :rtype: dict
        """
        return {}

    def getContainerPorts(self):
        """
        Ports to bind inside the container.

        :return: Dictionary of ports to bind
        :rtype: dict
        """
        return {}

    def getContainerVolumes(self):
        """
        Dictionary of volumes to mount inside container.

        :return: Dictionary of volumes
        :rtype: dict
        """
        return {}

    def getContainerWorkingDirectory(self):
        """
        Get working directory inside container.

        :return: Working directory, or None for default
        :rtype: str
        """
        return None

    def getNetwork(self):
        """
        Get Docker network to use with service container.

        :return: Docker network
        :rtype: docker.client.networks.Network
        """
        networkName = self.getNetworkName()
        try:
            return self.docker.networks.get(
                networkName
            )
        except docker.errors.NotFound:
            pass

        labels = Container.getLabels(self)
        labels["%s.project" % self.LABEL_PREFIX] = json.dumps(self.project)
        labels.pop(
            "%s.name" % self.LABEL_PREFIX
        )
        return self.docker.networks.create(
            networkName,
            labels=labels
        )

    def getContainer(self):
        """
        Get service Docker container if it exists.

        :return: Docker container
        :rtype: docker.client.containers.Container
        """
        if self._container:
            return self._container
        try:
            self._container = self.docker.containers.get(
                self.getContainerName()
            )
            return self._container
        except docker.errors.NotFound:
            pass
        return None

    def _createVolume(self, volumeId):
        """
        Create docker volume if it does not exist.
        """
        try:
            self.docker.volumes.get(volumeId)
        except docker.errors.NotFound:
            pass
        labels = Container.getLabels(self)
        self.docker.volumes.create(
            volumeId,
            labels=labels
        )

    def isRunning(self):
        """
        Determine if service container is currently running.

        :return: True if running
        :rtype: bool
        """
        container = self.getContainer()
        return container and container.status == "running"

    def getContainerIpAddress(self):
        """
        Get container local IP address.

        :return: Container IP address
        :rtype: str
        """
        container = self.getContainer()
        if not container:
            return ""
        return str(
            container.attrs.get(
                "NetworkSettings", {}
            ).get(
                "Networks", {}
            ).get(
                self.getNetworkName(), {}
            ).get(
                "IPAddress", ""
            )
        ).strip()

    def runCommand(self, command, user="root", shell="sh"):
        """
        Run a command inside container.

        :param command: Command to run
        :param user: User to run command as
        :param shell: Shell to use when running command
        :return: Output of command
        :rtype: str
        """
        container = self.getContainer()
        if not container:
            raise StateError(
                "Unable to get container '%s.'" % (
                    self.getContainerName()
                )
            )
        if not self.isRunning():
            raise StateError(
                "Container '%s' is not running." % self.getContainerName()
            )
        (exitCode, output) = container.exec_run(
            [
                shell, "-c", command
            ],
            user=user
        )
        if exitCode:
            self.logger.error(
                """
                Command execution error on container '%s'...
                (CMD=%s, USER=%s, EXIT_CODE=%s, OUTPUT=%s)
                """ % (
                    self.getContainerName(),
                    command,
                    user,
                    exitCode,
                    output.decode("utf-8")
                )
            )
            raise ContainerCommandError(
                "Command on container '%s' failed with exit code '%s.'" % (
                    self.getContainerName(),
                    exitCode
                )
            )
        return output.decode("utf-8")

    def uploadFile(self, uploadObj, path):
        """
        Upload a file to container.

        :param uploadObj: File object with data to push to container
        :param path: Path inside container to push data to
        """
        container = self.getContainer()
        if not container:
            raise StateError(
                "Unable to get container '%s.'" % (
                    self.getContainerName()
                )
            )
        if not self.isRunning():
            raise StateError(
                "Container '%s' is not running." % self.getContainerName()
            )
        with tempfile.NamedTemporaryFile(delete=True) as tarFile:
            with tarfile.open(fileobj=tarFile, mode="w") as tar:
                tarFileInfo = tarfile.TarInfo(
                    name=os.path.basename(path)
                )
                uploadObj.seek(0, io.SEEK_END)
                tarFileInfo.mtime = time.time()
                tarFileInfo.size = uploadObj.tell()
                uploadObj.seek(0)
                tar.addfile(
                    tarFileInfo,
                    uploadObj
                )
                
            tarFile.seek(0)
            container.put_archive(
                os.path.dirname(path),
                data=tarFile
            )

    def pullImage(self):
        """
        Pull image for container.
        """
        self.logger.info(
            "Pull '%s' image for '%s' container.",
            self.getBaseImage(),
            self.getContainerName()
        )
        self.docker.images.pull(self.getBaseImage())

    def getLabels(self):
        """
        Retrieve a list of labels to apply to container.
        """
        return {
            self.LABEL_PREFIX: "",
            "%s.project-uid" % self.LABEL_PREFIX: self.project.get("uid"),
            "%s.project-short-uid" % self.LABEL_PREFIX:
                self.project.get("short_uid"),
            "%s.name" % self.LABEL_PREFIX: self.getName()
        }

    def start(self):
        """
        Start Docker container for service.
        """
        self.logger.info(
            "Start '%s' container.",
            self.getContainerName()
        )
        self._container = None
        container = self.getContainer()
        # create docker container
        if not container:
            # add 'sys_admin' capability if mount volumes are to be used
            useMountVolumes = self.project.get("config", {}).get(
                "option_use_mount_volumes"
            )
            capAdd = []
            if useMountVolumes:
                capAdd.append("SYS_ADMIN")
            # create a network for docker container if not already created
            self.getNetwork()  
            # check if docker image name is defined
            # if not throw an error
            dockerImageName = self.getDockerImage()
            if not dockerImageName:
                raise StateError(
                    "No Docker image found for container '%s.'" % (
                        self.getContainerName()
                    )
                )
            # attempt to load docker image
            # pull the image if not found
            try:
                self.docker.images.get(dockerImageName)
            except docker.errors.ImageNotFound:
                self.pullImage()
            # create volumes if they don't exist
            volumes = self.getContainerVolumes()
            for volumeKey in volumes:
                if os.path.exists(volumeKey): continue
                self._createVolume(volumeKey)

            # create a docker container
            container = self.docker.containers.create(
                dockerImageName,
                name=self.getContainerName(),
                command=self.getContainerCommand(),
                detach=True,
                stdin_open=True,
                tty=True,
                environment=self.getContainerEnvironmentVariables(),
                extra_hosts=self.getContainerHosts(),
                network=self.getNetworkName(),
                ports=self.getContainerPorts(),
                volumes=volumes,
                working_dir=self.getContainerWorkingDirectory(),
                hostname=self.getContainerName(),
                privileged=bool(useMountVolumes),  # needed to mount
                cap_add=capAdd,
                labels=self.getLabels()
            )
        # if container already started do nothing
        if container.status == "running":
            return
        # start container
        container.start()
        self._container = None

    def stop(self):
        """
        Stop Docker container for service.
        """
        container = self.getContainer()
        if not container:
            return
        self.logger.info(
            "Stop '%s' container.",
            self.getContainerName()
        )
        container.stop()
        container.wait()
        container.remove()
        self._container = None

    def restart(self):
        """
        Restart Docker container for service.
        """
        self.stop()
        self.start()

    def commit(self):
        """
        Commit container in current state and create commit
        image.
        """
        container = self.getContainer()
        commitImage = self.getCommitImage().split(":")
        container.commit(
            commitImage[0],
            commitImage[1]
        )
        self._hasCommitImage = None

    def purge(self, dry=False):
        """
        Delete image and volumes for this container.

        :param dry: Don't perform purge, only list items to purge
        """
        # stop container if running
        if self.isRunning():
            self.stop()
        # delete committed image
        if self.getDockerImage() == self.getCommitImage():
            if not dry:
                self.docker.images.remove(self.getCommitImage())
            self.logger.info(
                "Delete '%s' Docker image.",
                self.getCommitImage()
            )
        # delete volumes
        for volumeName in self.getContainerVolumes():
            try:
                volume = self.docker.volumes.get(volumeName)
                if not dry:
                    volume.remove()
                self.logger.info(
                    "Deleted '%s' Docker volume.",
                    self.getVolumeName()
                )
            except docker.errors.NotFound:
                pass

    def shell(self, cmd="bash", user="root", stdin=None):
        """
        Create an interactive shell inside container.

        :param cmd: Command to run
        :param user: User to run as
        :param stdin: Stdin file object
        """
        if not self.isRunning():
            raise StateError(
                "Container '%s' is not running." % self.getContainerName()
            )
        container = self.getContainer()

        # has stdin
        if stdin and not stdin.isatty():
            with tempfile.NamedTemporaryFile(delete=True) as stdinFile:
                stdinFile.write(stdin.buffer.raw.read())
                # upload stdin object to container
                self.uploadFile(
                    stdinFile,
                    "/stdin.txt"
                )
            # inject stdin in to original command
            cmd = ["sh", "-c", "cat /stdin.txt | %s && rm /stdin.txt" % cmd]
            (_, output) = container.exec_run(
                cmd,
                user = "root"
            )
            # log output
            outputStr = output.decode("utf-8")
            if outputStr:
                self.logger.info(outputStr)
            return

        execId = self.docker.api.exec_create(
            container.id,
            cmd,
            tty=True,
            stdin=True,
            user=user
        )
        operation = ExecOperation(self.docker.api, execId)
        PseudoTerminal(self.docker.api, operation).start()
