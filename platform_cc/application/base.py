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
import json
import base64
import docker
import logging
import io
import random
import string
import collections
import platform
from nginx.config.api import Location
from nginx.config.api.options import KeyValueOption
from nginx.config.api.options import KeyValuesMultiLines
from nginx.config.api.options import KeyOption
from ..core.container import Container
from ..parser.routes import RoutesParser
from ..exception.state_error import StateError
from ..exception.container_command_error import ContainerCommandError
from ..services.base import BasePlatformService


class BasePlatformApplication(Container):

    """
    Base class for managing Platform.sh applications.
    """

    """
    Directory inside container to mount application to.
    """
    APPLICATION_DIRECTORY = "/app"

    """
    Directory inside container to mount storage to.
    """
    STORAGE_DIRECTORY = "/mnt/storage"

    """ Port to use for TCP upstream. """
    TCP_PORT = 8001

    """ Socket path to use for upstream. """
    SOCKET_PATH = "/run/app.sock"

    def __init__(self, project, config):
        """
        Constructor.

        :param project: Project data
        :param config: Service configuration
        """
        self.config = collections.OrderedDict(config)
        Container.__init__(
            self,
            project,
            self.config.get(
                "name",
                ""
            )
        )
        self.logger = logging.getLogger(
            "%s.%s.%s" % (
                __name__,
                self.project.get("short_uid"),
                self.getName()
            )
        )

    def isOSX(self):
        return platform.system() == "Darwin"

    def getContainerVolumes(self):
        useNFSVolumes = self.project.get("config", {}).get(
            "option_use_nfs_volumes"
        )
        if useNFSVolumes and self.isOSX():
            return {
                os.path.abspath(
                    self.config.get(
                        "_path", self.project.get("path"))
                    ).strip('/').replace("/","-"): {
                    "bind": self.APPLICATION_DIRECTORY,
                    "mode": "rw"
                },
                self.getVolumeName(): {
                    "bind": self.STORAGE_DIRECTORY,
                    "mode": "rw"
                }
            }
        else:
            return {
                os.path.abspath(
                    self.config.get(
                        "_path", self.project.get("path"))
                    ): {
                    "bind": self.APPLICATION_DIRECTORY,
                    "mode": "rw"
                },
                self.getVolumeName(): {
                    "bind": self.STORAGE_DIRECTORY,
                    "mode": "rw"
                }
            }

    def getApplicationVariables(self):
        output = {}
        varSources = [
            self.config.get("variables", {}),
            self.project.get("variables", {})
        ]
        for source in varSources:
            for _, k in enumerate(source):
                v = source[k]
                if type(v) is dict or type(v) is list:
                    for _, sk in enumerate(v):
                        output["%s:%s" % (k, sk)] = v[sk]
                    continue
                output[k] = v
        return output

    def getContainerEnvironmentVariables(self):
        # get platform relationships
        platformRelationships = {}
        for key, value in self.config.get("relationships", {}).items():
            value = value.strip().split(":")
            platformRelationships[key] = [
                self.project.get("services", {})
                    .get(value[0], {})
                    .get("platform_relationships", {})
                    .get(value[1])
            ]
        routesParser = RoutesParser(self.project)
        # get subnet from project network, used for trusted proxy
        network = self.getNetwork()
        trustedProxies = "%s,127.0.0.1" % (
            str(
                network.attrs.get("IPAM", {}).get(
                    "Config", [{}]
                )[0].get("Subnet")
            )
        )
        try:
            bridgeNetwork = self.docker.networks.get("bridge")
            trustedProxies = "%s,%s" % (
                str(
                    bridgeNetwork.attrs.get("IPAM", {}).get(
                        "Config", [{}]
                    )[0].get("Subnet")
                ),
                trustedProxies
            )
        except docker.errors.NotFound:
            pass
        # set env vars
        envVars = {
            "PLATFORM_APP_DIR":           self.APPLICATION_DIRECTORY,
            "PLATFORM_APPLICATION":       "",
            "PLATFORM_APPLICATION_NAME":  self.getName(),
            "PLATFORM_BRANCH":            "",
            "PLATFORM_DOCUMENT_ROOT":     "/",
            "PLATFORM_ENVIRONMENT":       "platform_cc",
            "PLATFORM_PROJECT":           self.project.get("uid", ""),
            "PLATFORM_RELATIONSHIPS":     base64.b64encode(
                bytes(str(json.dumps(platformRelationships)).encode("utf-8"))
            ).decode("utf-8"),
            "PLATFORM_ROUTES":            base64.b64encode(
                bytes(str(
                    json.dumps(routesParser.getRoutesEnvironmentVariable())
                ).encode("utf-8"))
            ).decode("utf-8"),
            "PLATFORM_TREE_ID":           ''.join(random.choice(
                string.ascii_lowercase + string.digits) for _ in range(40)
            ),
            "PLATFORM_VARIABLES":         base64.b64encode(
                bytes(
                    str(
                        json.dumps(self.getApplicationVariables())
                    ).encode("utf-8")
                )
            ).decode("utf-8"),
            "PLATFORM_PROJECT_ENTROPY"  : self.project.get("entropy", ""),
            "TRUSTED_PROXIES"           : trustedProxies,
            "SYMFONY_TRUSTED_PROXIES"   : trustedProxies,
            "PORT"                      : self.TCP_PORT,
            "SOCKET"                    : self.SOCKET_PATH
        }
        # set env vars from app variables
        for key, value in self.getApplicationVariables().items():
            if not key.startswith("env:"): continue
            key = key[4:]
            envVars[key.strip().upper()] = str(value)
        return envVars

    def getContainerWorkingDirectory(self):
        return self.APPLICATION_DIRECTORY

    def getType(self):
        """
        Get application type.

        :return: Application type
        :rtype: str
        """
        return self.config.get("type")

    def _generateNginxPassthruOptions(self, locationConfig={}):
        """
        Get options to generate nginx passthru.

        :param locationConfig: Dict containing location configuration
        :return: List of nginx block values
        :rtype: list
        """
        upstreamConf = self.config.get("web", {}).get(
            "upstream", {"socket_family": "tcp", "protocol": "http"}
        )
        output = []
        # tcp port, proxy pass
        if (
            upstreamConf.get("socket_family") == "tcp" and
            upstreamConf.get("protocol") == "http"
        ):
            output.append(
                KeyValueOption(
                    "proxy_pass",
                    "http://127.0.0.1:%d" % self.TCP_PORT
                )
            )
            output.append(
                KeyValueOption("proxy_set_header", "Host $host")
            )
        # tcp port, fastcgi
        elif (
            upstreamConf.get("socket_family") == "tcp" and
            upstreamConf.get("protocol") == "fastcgi"
        ):
            output.append(
                KeyValueOption("fastcgi_pass", "127.0.0.1:%d" % self.TCP_PORT)
            )
            output.append(
                KeyValueOption("include", "fastcgi_params")
            )
            output.append(
                KeyValueOption("set", "$path_info $fastcgi_path_info")
            )
        # socket, proxy pass
        elif (
            upstreamConf.get("socket_family") == "socket" and
            upstreamConf.get("protocol") == "http"
        ):
            output.append(
                KeyValueOption("proxy_pass", "unix:%s" % self.SOCKET_PATH)
            )
            output.append(
                KeyValueOption("proxy_set_header", "Host $host")
            )
        # socket, fastcgi
        elif (
            upstreamConf.get("socket_family") == "socket" and
            upstreamConf.get("protocol") == "fastcgi"
        ):
            output.append(
                KeyValueOption("fastcgi_pass", "unix:%s" % self.SOCKET_PATH)
            )
            output.append(
                KeyValueOption("include", "fastcgi_params")
            )
            output.append(
                KeyValueOption("set", "$path_info $fastcgi_path_info")
            )
        return output

    def _generateNginxLocationOptions(self, locationConfig={}):
        """
        Get options to generate Nginx location.
        """
        output = []
        # alias
        root = locationConfig.get("root", "") or ""
        output.append(
            KeyValueOption(
                "alias",
                "%s/" % ("%s/%s" % (
                    self.APPLICATION_DIRECTORY, root.strip("/")
                )).rstrip("/")
            )
        )
        # index
        index = locationConfig.get("index", [])
        if type(index) is not list:
            index = [index]
        if index:
            output.append(
                KeyValueOption("index", " ".join(index))
            )
        # expires
        expires = locationConfig.get("expires", "-1s")
        if type(expires) is not str:
            expires = "%ds" % expires
        if locationConfig.get("allow", False):
            output.append(
                KeyValueOption("expires", expires)
            )
            if expires != "-1s":
                output.append(
                    KeyValueOption("gzip", "on")
                )
                output.append(
                    KeyValueOption("brotli", "on")
                )
        # headers
        headers = locationConfig.get("headers", {})
        if headers:
            output.append(
                KeyValuesMultiLines(
                    "add_header",
                    ["\"%s\" \"%s\"" % (k, v) for k, v in headers.items()]
                )
            )
        return output

    def _generateNginxRootLocation(self, path, locationConfig={}):
        """
        Generate nginx root location configuration.
        """
        locationConfig = locationConfig.copy()
        locationConfig["rules"] = {}
        rootLocation = Location(
            "= \"%s\"" % path.rstrip("/"),
            *self._generateNginxLocationOptions(locationConfig)
        )
        return rootLocation

    def _generateNginxLocation(self, path, locationConfig={}):
        """
        Generate nginx location configuration for given path.

        :param path: Location path
        :param locationConfig: Dict of location configuration
        :return: Nginx location
        :rtype: Location
        """
        # create location
        pathStrip = "/%s/" % path.strip("/")
        if pathStrip == "//":
            pathStrip = "/"
        location = Location(
            "\"%s\"" % pathStrip,
            *self._generateNginxLocationOptions(locationConfig)
        )
        # passthru
        passthru = locationConfig.get("passthru", False)
        if passthru:
            passthruLocation = Location(
                "~ /",
                allow="all",
                *self._generateNginxPassthruOptions(locationConfig)
            )
            location.sections.add(passthruLocation)
        # rules, add extra locations
        for key, val in locationConfig.get("rules", {}).items():
            ruleConfig = locationConfig.copy()
            ruleConfig.update(val)
            ruleConfig["rules"] = {}
            location.sections.add(
                self._generateNginxLocation(
                    key,
                    ruleConfig
                )
            )
        # output
        return location

    def generateNginxConfig(self):
        """
        Generate configuration for nginx specific to application.

        :return: Nginx configuration
        :rtype: str
        """
        self.logger.info(
            "Generate application Nginx configuration."
        )
        locations = self.config.get("web", {}).get("locations", {})
        if not locations or len(locations) == 0:
            locations["/"] = {
                "allow":      False,
                "passthru":   True
            }
        output = "charset UTF-8;\n"
        for path, config in locations.items():
            output += str(self._generateNginxRootLocation(
                path, config
            ))
            output += str(self._generateNginxLocation(
                path, config
            ))
        return output

    def getMounts(self):
        """
        Get mount points.
        """
        configMounts = self.config.get("mounts", {})
        output = {}
        for mountDest, config in configMounts.items():
            mountSrc = ""
            if type(config) is dict:
                mountSrc = config.get("source_path", "").strip("/")
            elif type(config) is str:
                localMountPrefx = "shared:files/"
                if not config.startswith(localMountPrefx):
                    continue
                mountSrc = config[len(localMountPrefx):].strip("/")
            else:
                continue
            output[mountSrc.strip("/")] = mountDest.strip("/")
        return output

    def setupMounts(self):
        """
        Setup application defined mount points.
        """
        # project option 'use_mount_volumes' must be true
        mounts = self.getMounts()
        self.logger.info(
            "Found %s mount point(s).",
            len(mounts)
        )
        for mountSrc, mountDest in mounts.items():
            self.logger.debug(
                "Bind mount point '%s.'.",
                mountSrc
            )
            mountSrc = os.path.join(
                self.STORAGE_DIRECTORY,
                mountSrc.strip("/")
            )
            mountDest = os.path.join(
                self.APPLICATION_DIRECTORY,
                mountDest.strip("/")
            )
            self.runCommand(
                "mkdir -p %s && mkdir -p %s" % (
                    mountSrc,
                    mountDest,
                ),
                "root"
            )
            if not self.project.get("config", {}).get(
                "option_use_mount_volumes"
            ):
                continue
            self.runCommand(
                "mount -o user_xattr --bind %s %s" % (
                    mountSrc,
                    mountDest
                ),
                "root"
            )

    def installSsh(self):
        """
        Install SSH key and known hosts file.
        """
        sshDatas = [
            ["ssh_key", "/app/.ssh/id_rsa"],
            ["ssh_known_hosts", "/app/.ssh/known_hosts"]
        ]
        try:
            self.runCommand(
                "mkdir -p /app/.ssh && chown -f -R web:web /app/.ssh"
            )
        except ContainerCommandError:
            pass
        for sshData in sshDatas:
            data = self.project.get("config", {}).get(sshData[0])
            if not data:
                continue
            self.logger.info(
                "Install '%s' from project config." % sshData[0]
            )
            data = base64.b64decode(data)
            dataFileObject = io.BytesIO(data)
            self.uploadFile(
                dataFileObject,
                "/tmp/.ssh_file"  # upload to tmp and copy
            )
            try:
                self.runCommand(
                    "mv /tmp/.ssh_file %s && chmod -f 0600 %s" % (
                        sshData[1],
                        sshData[1]
                    )
                )
            except ContainerCommandError:
                pass

    def installCron(self):
        """
        Install cron tasks and enable cron in application container.
        """
        # cron must be enabled via options
        if not self.project.get("config", {}).get("option_enable_cron"):
            return
        # create cron directory if not exist
        self.runCommand(
            "mkdir -p /etc/cron.d"
        )
        # itterate crons make cron files
        crons = self.config.get("crons", {})
        self.logger.info(
            "Installing %s cron task(s)." % str(len(crons))
        )
        for name, cron in crons.items():
            spec = cron.get("spec", "*/5 * * * *")  # default is every 5m
            cmd = cron.get("cmd", "")
            if not cmd:
                continue
            self.logger.debug(
                "Installing '%s' cron." % name
            )
            fileObj = io.BytesIO(
                bytes(
                    "%s web %s" % (
                        spec,
                        cmd
                    )
                )
            )
            self.uploadFile(
                fileObj,
                "/etc/cron.d/%s" % name
            )
        # start cron
        self.logger.info("Start cron daemon.")
        self.runCommand("cron")

    def prebuild(self):
        """
        Perform tasks on container prior to build process.
        """
        # delete committed image
        if self.getDockerImage() == self.getCommitImage():
            # stop container if running
            if self.isRunning():
                self.stop()
            self.docker.images.remove(self.getCommitImage())
            self.logger.info(
                "Delete '%s' Docker image.",
                self.getCommitImage()
            )
            self._hasCommitImage = False
        # start container
        if not self.isRunning():
            BasePlatformApplication.start(self, False, False)

    def build(self):
        """
        Run commands needed to get container ready for given
        application. Also runs build hooks commands.
        """
        self.prebuild()
        self.logger.info(
            "Building application."
        )
        # install ssh
        self.installSsh()
        # run build hooks
        print("- INTERACTIVE SHELL --------------------------------------")
        self.shell(
            """
            bash -c '%s'
            """ % self.config.get("hooks", {}).get("build", ""),
            "root"
        )
        print("----------------------------------------------------------")
        # commit container
        self.logger.info(
            "Commit container."
        )
        self.commit()
        self.stop()

    def deploy(self):
        """
        Run deploy hook commands.
        """
        self.logger.info(
            "Run deploy hooks."
        )
        print("- INTERACTIVE SHELL --------------------------------------")
        self.shell(
            """
            bash -c '%s'
            """ % self.config.get("hooks", {}).get("deploy", ""),
            "web"
        )
        print("----------------------------------------------------------")

    def getLabels(self):
        labels = Container.getLabels(self)
        labels["%s.config" % Container.LABEL_PREFIX] = json.dumps(self.config)
        labels["%s.type" % Container.LABEL_PREFIX] = "application"
        return labels

    def startServices(self):
        """ Start extra services ran in the app container. """
        # nginx
        self.logger.info(
            "Start Nginx."
        )
        nginxConfFileObj = io.BytesIO(
            bytes(str(self.generateNginxConfig()).encode("utf-8"))
        )
        self.uploadFile(
            nginxConfFileObj,
            "/usr/local/nginx/conf/app.conf"
        )
        self.runCommand(
            """
            /usr/local/nginx/sbin/nginx -s stop || true
            /usr/local/nginx/sbin/nginx
            """
        )
        # cron
        self.installCron()

    def start(self, requireServices=True, buildIfNotBuilt=True, setupMounts=True):
        # ensure all required services are available
        if requireServices:
            projectServices = self.project.get("services", {})
            serviceNames = list(self.config.get("relationships", {}).values())
            for serviceName in serviceNames:
                serviceName = serviceName.strip().split(":")[0]
                projectService = projectServices.get(serviceName)
                if not projectService or not projectService.get("running"):
                    if projectService and projectService.get("start_group") not in [BasePlatformService.START_PRE_APP_A, BasePlatformService.START_PRE_APP_B]:
                        continue
                    raise StateError(
                        "Application '%s' depends on service '%s' which is not running." % (
                            self.getName(),
                            serviceName
                        )
                    )
        # start container
        Container.start(self)
        container = self.getContainer()
        if not container:
            return
        # setup mount points
        if setupMounts: self.setupMounts()
        # not yet built/provisioned
        if self.getDockerImage() == self.getBaseImage() and buildIfNotBuilt:
            self.build()
            return self.start(requireServices)
