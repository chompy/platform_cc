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
from platform_cc.container import Container

class BasePlatformService(Container):
    """
    Base class for Platform.sh services.
    """

    def __init__(self, project, config, dockerClient = None):
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
            "running"                   : self.isRunning(),
            "ip"                        : self.getContainerIpAddress(),
            "platform_relationships"    : {}
        }

    def getLabels(self):
        labels = Container.getLabels(self)
        labels["%s.config" % Container.LABEL_PREFIX] = json.dumps(self.config)
        labels["%s.type" % Container.LABEL_PREFIX] = "service"
        return labels

    def start(self):
        self.logger.info("Start '%s' service." % self.getName())
        # if not platform.sh compatiable service and service definition
        # is in main service.yaml file warn user to consider
        # moving service definition to service.pcc.yaml
        if not self.isPlatformShCompatible():
            # TODO
            pass

        Container.start(self)