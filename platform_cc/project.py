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

from __future__ import absolute_import
import os
import json
import time
import hashlib
import base36
import random
import io
import docker
import logging
from platform_cc.variables import getVariableStorage
from platform_cc.variables.var_json import JsonVariables
from platform_cc.parser.services import ServicesParser
from platform_cc.parser.applications import ApplicationsParser
from platform_cc.parser.routes import RoutesParser
from platform_cc.services import getService
from platform_cc.application import getApplication
from platform_cc.router import PlatformRouter
from platform_cc.exception.state_error import StateError

class PlatformProject:
    """
    Container class for all elements of a Platform.sh
    project. (Applications, services, variables, etc).
    """

    """ Filename of project config file. """
    PROJECT_CONFIG_FILE = ".pcc_project.json"

    """ Salt used to generate project unique ids. """
    HASH_SALT = "6fabb8b0ee9&(2cae2eb26306cdc51012f180eb$NBd!a0e"

    def __init__(self, path):
        """
        Constructor.

        :param path: Path to project root
        """

        # get logger
        self.logger = logging.getLogger(__name__)

        # set project path
        self.path = str(path)

        # validate project path
        if not os.path.isdir(self.path):
            raise ValueError("Invalid project path.")

        # load config (use JsonVariables class to do this
        # as it already contains the functionality)
        self.config = JsonVariables(
            self.path,
            {
                "variables_json_filename" : self.PROJECT_CONFIG_FILE
            }
        )

        # generate uid if it does not exist
        if not self.config.get("uid"):
            self.config.set("uid", self._generateUid())

        # define container 'web' user id for nginx
        # and application process to use
        # this should be the same user id as the current
        # user so that permissions in and out of the container match
        if not self.config.get("web_user_id") or self.config.get("web_user_id") <= 0:
            currentUserId = os.getuid()
            if currentUserId <= 0:
                currentUserId = 1000
            self.config.set("web_user_id", currentUserId)

        # get variable storage
        self.variables = getVariableStorage(
            self.path,
            self.config
        )

        # define services parser
        self._servicesParser = None

        # list of initialized services
        self._services = []

        # define applications parser
        self._applicationsParser = None

        # list of initialized applications
        self._applications = []

        # define routes parser
        self._routesParser = None

        # define router
        self._router = None

    def _generateUid(self):
        """
        Generate a unique id for the project.

        :param path: Path to project root
        :return: Unique id string
        :rtype: str
        """
        return base36.dumps(
            int(
                hashlib.sha256(
                    (
                        "%s-%s-%s-%s" % (
                            self.HASH_SALT,
                            self.path,
                            str(random.random()),
                            str(time.time())
                        )
                    ).encode("utf-8")
                ).hexdigest(),
                16
            )
        )

    def getUid(self):
        """ 
        Get project unique id.

        :return: Project unique id
        :rtype: str
        """
        return self.config.get("uid")

    def getShortUid(self):
        """
        Get shortened project unique id.

        :return: Shortened project unique id
        :rtype: str
        """
        return self.getUid()[0:6]

    def getEntropy(self):
        """
        Get entropy value for project.

        :return: Unique random string
        :rtype: str
        """
        entropy = self.config.get("entropy")
        if entropy: return entropy
        entropy = base36.dumps(
            int(
                hashlib.sha256(
                    (
                        "%s-%s-%s-%s" % (
                            self.getUid(),
                            str(random.random()),
                            str(time.time()),
                            self.HASH_SALT
                        )
                    ).encode("utf-8")
                ).hexdigest(),
                16
            )
        )
        self.config.set("entropy", entropy)
        return entropy

    def getProjectData(self):
        """
        Get dictionary containing project data
        needed by other child objects.

        :return: Dictionary of project data
        :rtype: dict
        """
        serviceData = {}
        for service in self._services:
            if not service: continue
            serviceData[service.getName()] = service.getServiceData()
        return {
            "path"              : self.path,
            "uid"               : self.getUid(),
            "short_uid"         : self.getShortUid(),
            "entropy"           : self.getEntropy(),
            "config"            : self.config.all(),
            "variables"         : self.variables.all(),
            "services"          : serviceData # service data needed by applications
        }

    def getServicesParser(self):
        """
        Get services parser.

        :return: Services parser
        :rtype: .parser.services.ServicesParser
        """
        if self._servicesParser:
            return self._servicesParser
        self.logger.debug(
            "Build services parser for project '%s.'",
            self.getShortUid()
        )
        self._servicesParser = ServicesParser(self.path)
        return self._servicesParser

    def getApplicationsParser(self):
        """
        Get applications parser.

        :return: Applications parser
        :rtype: .parser.applications.ApplicationsParser
        """
        if self._applicationsParser:
            return self._applicationsParser
        self.logger.debug(
            "Build applications parser for project '%s.'",
            self.getShortUid()
        )
        self._applicationsParser = ApplicationsParser(self.path)
        return self._applicationsParser

    def getRoutesParser(self):
        """
        Get routes parser.

        :return: Routes parser
        :rtype: .parser.routes.RoutesParser
        """
        if self._routesParser:
            return self._routesParser
        self.logger.debug(
            "Build routes parser for project '%s.'",
            self.getShortUid()
        )            
        self._routesParser = RoutesParser(self.getProjectData())
        return self._routesParser

    def getService(self, name):
        """
        Get a service handler.

        :param name: Service name
        :return: Service handler
        :rtype: .service.base.BasePlatformService
        """
        servicesParser = self.getServicesParser()
        name = str(name)
        for service in self._services:
            if service and service.getName() == name:
                return service
        serviceConfig = servicesParser.getServiceConfiguration(name)
        self.logger.debug(
            "Build service handler '%s' for project '%s.'",
            name,
            self.getShortUid()
        )
        service = getService(
            self.getProjectData(),
            serviceConfig
        )
        self._services.append(service)
        return service

    def getApplication(self, name):
        """
        Get an application handler.

        :param name: Application name
        :return: Application handler
        :rtype: .application.base.BasePlatformApplication
        """
        applicationsParser = self.getApplicationsParser()
        name = str(name)
        for application in self._applications:
            if application and application.getName() == name:
                return application       
        appConfig = applicationsParser.getApplicationConfiguration(name)
        # init all service dependencies for app
        appServiceDependencies = list(appConfig.get("relationships", {}).values())
        for serivceName in appServiceDependencies:
            self.getService(serivceName.split(":")[0])
        # get app
        self.logger.debug(
            "Build application handler '%s' for project '%s.'",
            name,
            self.getShortUid()
        )        
        application = getApplication(
            self.getProjectData(),
            appConfig
        )
        self._applications.append(application)
        return application

    def getRouter(self):
        """
        Get main project router.
        """
        if self._router: return self._router
        self._router = PlatformRouter()
        return self._router

    def addRouter(self):
        """
        Add this project to the router.
        """
        self.logger.info(
            "Add project '%s' to router.",
            self.getShortUid()
        )
        # get router
        router = self.getRouter()
        if not router.isRunning(): router.start()
        # retrieve all applications
        appParser = self.getApplicationsParser()
        for appName in appParser.getApplicationNames():
            self.getApplication(appName)
        if len(self._applications) == 0:
            raise Exception("Project must contain at least one application.")
        # generate nginx config
        nginxConfig = router.generateNginxConfig(self._applications)
        # upload nginx config to router
        nginxConfigFile = io.BytesIO(
            bytes(str(nginxConfig).encode("utf-8"))
        )
        router.uploadFile(
            nginxConfigFile,
            os.path.join(
                router.NGINX_PROJECT_CONF_PATH,
                "%s.conf" % self.getShortUid()
            )
        )
        # add router to project network
        network = self._applications[0].getNetwork()
        try:
            network.connect(
                router.getContainer()
            )
        except docker.errors.APIError:
            pass
        # restart router
        router.getContainer().restart()

    def removeRouter(self):
        """
        Remove this project from the router.
        """
        self.logger.info(
            "Remove project '%s' from router.",
            self.getShortUid()
        )
        # get router
        router = self.getRouter()
        if not router.isRunning(): return
        # delete conf file in router
        router.runCommand(
            "rm -f %s.conf" % (
                os.path.join(
                    router.NGINX_PROJECT_CONF_PATH,
                    self.getShortUid()
                )
            )
        )
        router.getContainer().restart()