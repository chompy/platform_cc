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

import logging
import json
import docker
from platform_cc.container import Container
from platform_cc.parser.services import ServicesParser


class BasePlatformService(Container):
    """
    Base class for Platform.sh services.
    """

    START_PRE_APP_A = "preapp-a"
    START_PRE_APP_B = "preapp-b"
    START_POST_APP_A = "postapp-a"
    START_POST_APP_B = "postapp-b"

    def __init__(self, project, config, dockerClient=None):
        """
        Constructor.

        :param project: Project data
        :param config: Service configuration
        :param dockerClient: Docker client
        """
        self.config = dict(config)
        Container.__init__(
            self,
            project,
            self.config.get(
                "_name",
                self.config.get(
                    "_type"
                )
            ),
            dockerClient
        )
        self.logger = logging.getLogger(
            "%s.%s.%s" % (
                __name__,
                self.project.get("short_uid"),
                self.getName()
            )
        )
        isDefaultConf = self.config.get("_is_default_config", True)
        if isDefaultConf and not self.isPlatformShCompatible():
            self.logger.warn(
                """
                Service '%s/%s' is not compatible with Platform.sh but
                was defined in '%s', consider moving it to '%s.'
                """ % (
                    self.getName(),
                    self.getType(),
                    ServicesParser.YAML_PATHS[0],
                    ServicesParser.YAML_PATHS[1]
                )
            )

    def getStartGroup(self):
        """
        Define start group to determine when to start
        service.
        """
        return self.START_PRE_APP_A

    def getType(self):
        """
        Get service type.

        :return: Service type
        :rtype: str
        """
        return self.config.get(
            "_type"
        )

    def isPlatformShCompatible(self):
        """
        Whether or not this service is designed to
        be compatible with platform.sh.

        :rtype: bool
        """
        return True

    def getServiceData(self):
        """
        Get data needed to access service for use by applications.

        :return: Dictionary containing service data
        :rtype: dict
        """
        return {
            "running":                    self.isRunning(),
            "ip":                         self.getContainerIpAddress(),
            "platform_relationships":     {},
            "start_group":                self.getStartGroup()
        }

    def getLabels(self):
        labels = Container.getLabels(self)
        labels["%s.config" % Container.LABEL_PREFIX] = json.dumps(self.config)
        labels["%s.type" % Container.LABEL_PREFIX] = "service"
        return labels

    def getContainerCommand(self):
        # override default so that we can inject custom scripts
        try:
            image = self.docker.images.get(self.getDockerImage())
        except docker.errors.ImageNotFound:
            self.docker.images.pull(self.getDockerImage())
            image = self.docker.images.get(self.getDockerImage())
        return image.attrs.get("Config", {}).get("Cmd", [])

    def start(self):
        self.logger.info("Start '%s' service." % self.getName())
        # if not platform.sh compatiable service and service definition
        # is in main service.yaml file warn user to consider
        # moving service definition to service.pcc.yaml
        if not self.isPlatformShCompatible():
            # TODO
            pass
        Container.start(self)
