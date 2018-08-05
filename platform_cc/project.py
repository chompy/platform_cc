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
from platform_cc.variables.base import BasePlatformVariables
from platform_cc.parser.services import ServicesParser
from platform_cc.parser.applications import ApplicationsParser
from platform_cc.parser.routes import RoutesParser
from platform_cc.services import getService
from platform_cc.application import getApplication
from platform_cc.router import PlatformRouter
from platform_cc.container import Container
from platform_cc.exception.state_error import StateError
from platform_cc.exception.parser_error import ParserError

class PlatformProject:
    """
    Container class for all elements of a Platform.sh
    project. (Applications, services, variables, etc).
    """

    """ Filename of project config file. """
    PROJECT_CONFIG_FILE = ".pcc_project.json"

    """ Filename to use when storing variables. """
    PROJECT_VAR_STORAGE_FILE = ".pcc_variables.json"

    """ Salt used to generate project unique ids. """
    HASH_SALT = "6fabb8b0ee9&(2cae2eb26306cdc51012f180eb$NBd!a0e"

    def __init__(self, projectConfig, projectVars, projectPath = None):
        """
        Constructor.

        :param projectConfig: Project configuration
        :param projectVars: Project variables
        :param projectPath: Path to project
        """

        # type check
        if not isinstance(projectConfig, BasePlatformVariables):
            raise ValueError("Project config should be an instance of 'BasePlatformVariables.'")
        if not isinstance(projectVars, BasePlatformVariables):
            raise ValueError("Project variables should be an instance of 'BasePlatformVariables.'")

        # set project path, enforce string
        self.path = None
        if projectPath:
            self.path = str(projectPath)

        # validate project path
        print(self.path)
        if self.path and not os.path.isdir(self.path):
            raise ValueError("Project path does not exist.")

        # get logger
        self.logger = logging.getLogger(__name__)

        # set project config
        self.config = projectConfig

        # set project variables
        self.variables = projectVars

        # generate uid if it does not exist
        if not self.config.get("uid"):
            self.config.set("uid", self._generateUid())

        # define container 'web' user id for nginx
        # and application process to use
        # this should be the same user id as the current
        # user so that permissions in and out of the container match
        if not self.config.get("web_user_id") or int(self.config.get("web_user_id")) <= 0:
            try:
                currentUserId = os.getuid()
            except AttributeError:
                currentUserId = 1000
            if currentUserId <= 0:
                currentUserId = 1000
            self.config.set("web_user_id", currentUserId)

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

    @staticmethod
    def fromPath(projectPath):
        """
        Load instance of project from path to project root.

        :param projectPath: Path to project
        :rtype: PlatformProject
        """

        # enforce string
        projectPath = str(projectPath)

        # validate project path
        if not os.path.isdir(projectPath):
            raise ValueError("Project path does not exist.")

        # load config (use JsonFileVariables class to do this
        # as it already contains the functionality)
        projectConfig = getVariableStorage(
            {
                "storage_handler"        : "json_file",
                "json_path"              : os.path.join(projectPath, PlatformProject.PROJECT_CONFIG_FILE)
            }
        )

        # get variable storage
        projectVars = getVariableStorage(
            {
                "storage_handler"        : "json_file",
                "json_path"              : os.path.join(projectPath, PlatformProject.PROJECT_VAR_STORAGE_FILE)
            }
        )

        # create project object
        return PlatformProject(
            projectConfig,
            projectVars,
            projectPath
        )

    @staticmethod
    def getDockerClient():
        """
        Get the Docker client.
        """
        return docker.from_env(
            timeout = 300 # 5 minutes
        )

    @staticmethod
    def fromDocker(projectUid):
        """
        Load instance of project from values stored in
        Docker labels using the project uid.

        :param projectUid: Project unique id
        :rtype: PlatformProject
        """

        # fetch project network
        dockerClient = PlatformProject.getDockerClient()
        networkList = dockerClient.networks.list(
            filters = {
                "label" : ["%s.project-uid=%s" % (
                    Container.LABEL_PREFIX, projectUid
                )]
            }
        )
        if not networkList:
            networkList = dockerClient.networks.list(
                filters = {
                    "label" : ["%s.project-short-uid=%s" % (
                        Container.LABEL_PREFIX, projectUid
                    )]
                }
            )
        if not networkList:
            raise ValueError("Cannot find active project with uid '%s' in current Docker environment." % projectUid)

        # fetch network labels
        labels = networkList[0].attrs.get("Labels", {})

        # get project data
        projectData = json.loads(labels.get("%s.project" % Container.LABEL_PREFIX))

        # get config from project data
        projectConfig = getVariableStorage(
            {
                "storage_handler"           : "dict",
                "dict_vars"                 : projectData.get("config", {})
            }
        )

        # get vars from project data
        projectVars = getVariableStorage(
            {
                "storage_handler"           : "dict",
                "dict_vars"                 : projectData.get("variables", {})
            }
        )

        # create project object
        return PlatformProject(
            projectConfig,
            projectVars,
            None
        )

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
        if not self.path:
            raise StateError("Cannot get services parse without project path.")
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
        if not self.path:
            raise StateError("Cannot get applications parser without project path.")
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
        if not self.path:
            raise StateError("Cannot get routes parser without project path.")
        self.logger.debug(
            "Build routes parser for project '%s.'",
            self.getShortUid()
        )            
        self._routesParser = RoutesParser(self.getProjectData())
        return self._routesParser

    def getService(self, name = None):
        """
        Get a service handler.

        :param name: Service name, if not provided first found service is returned
        :return: Service handler
        :rtype: .service.base.BasePlatformService
        """

        # service already loaded
        for service in self._services:
            if service and (not name or service.getName() == name):
                return service

        # check docker
        dockerClient = PlatformProject.getDockerClient()
        dockerLabelFilters = [
            "%s.type=service" % (Container.LABEL_PREFIX),
            "%s.project-uid=%s" % (Container.LABEL_PREFIX, self.getUid())
        ]
        if name:
            dockerLabelFilters.append(
                "%s.name=%s" % (Container.LABEL_PREFIX, name)
            )
        matchedServices = dockerClient.containers.list(
            filters = {
                "label" : dockerLabelFilters
            }
        )
        if matchedServices:
            self.logger.debug(
                "Build service handler '%s' for project '%s' from Docker environment.",
                matchedServices[0].attrs.get("Config", {}).get("Labels", {}).get("%s.name" % Container.LABEL_PREFIX),
                self.getShortUid()
            )
            serviceConfig = json.loads(
                matchedServices[0].attrs.get("Config", {}).get("Labels", {}).get("%s.config" % Container.LABEL_PREFIX)
            )
            service = getService(
                self.getProjectData(),
                serviceConfig
            )
            self._services.append(service)
            return service

        # check path
        if self.path:
            servicesParser = self.getServicesParser()
            if not name:
                serviceNames = servicesParser.getServiceNames()
                if not serviceNames:
                    raise ParserError("No services defined.")
                name = servicesParser.getServiceNames()[0]
            name = str(name)
            self.logger.debug(
                "Build service handler '%s' for project '%s' from YAML configuration files.",
                name,
                self.getShortUid()
            )
            serviceConfig = servicesParser.getServiceConfiguration(name)
            service = getService(
                self.getProjectData(),
                serviceConfig
            )
            self._services.append(service)
            return service

        # not found
        raise StateError("Unable to find service with name '%s.'" % name)

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