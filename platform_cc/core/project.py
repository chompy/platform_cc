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
import collections
from ..variables import getVariableStorage
from ..variables.base import BasePlatformVariables
from ..parser.services import ServicesParser
from ..parser.applications import ApplicationsParser
from ..parser.routes import RoutesParser
from ..services import getService
from ..services.base import BasePlatformService
from ..application import getApplication
from .router import PlatformRouter
from .container import Container
from .config import PlatformConfig
from ..exception.state_error import StateError
from ..exception.parser_error import ParserError
from ..exception.container_not_found_error import ContainerNotFoundError
from ..exception.project_init_error import ProjectInitError

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
        if self.path and not os.path.isdir(self.path):
            raise ProjectInitError("Project path does not exist.")

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
            raise ProjectInitError("Project path does not exist.")

        # ensure a platform.sh related yaml file exists
        mustContainOneOf = ServicesParser.YAML_PATHS + RoutesParser.YAML_PATHS + ApplicationsParser.YAML_FILENAMES
        containsOne = False
        for i in mustContainOneOf:
            pathTo = os.path.join(projectPath, i)
            if os.path.exists(pathTo):
                containsOne = True
                break
        if not containsOne:
            raise ProjectInitError("No project configuration files found.")

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
                "json_path"              : os.path.join(projectPath, PlatformProject.PROJECT_VAR_STORAGE_FILE),
                "global_vars"            : PlatformConfig().getGlobalProjectVars()
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
            timeout = 1800 # 30 minutes
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
            raise ProjectInitError("Cannot find active project with uid '%s' in current Docker environment." % projectUid)

        # fetch network labels
        labels = networkList[0].attrs.get("Labels", {})

        # get project data
        projectData = json.loads(
            labels.get("%s.project" % Container.LABEL_PREFIX),
            object_pairs_hook=collections.OrderedDict
        )

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
            "services"          : serviceData, # service data needed by applications
            "_enable_service_routes" : self.config.get("option_enable_service_routes")
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
            raise StateError("Cannot get services parser without project path.")
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

    def dockerFetch(self, filterType = None, filterName = None, all = False):
        """
        Fetch applications and services from Docker environment.

        :param filterType: Filter type of container (application or service)
        :param filterName: Filter by container name
        :param all: If true all containers, even non running, are fetched
        :return: List containing containers
        """

        dockerClient = PlatformProject.getDockerClient()
        dockerLabelFilters = [
            "%s.project-uid=%s" % (Container.LABEL_PREFIX, self.getUid())
        ]
        if filterType:
            dockerLabelFilters.append(
                "%s.type=%s" % (Container.LABEL_PREFIX, filterType)
            )            
        if filterName:
            dockerLabelFilters.append(
                "%s.name=%s" % (Container.LABEL_PREFIX, filterName)
            )
        containerList = dockerClient.containers.list(
            filters = {
                "label" : dockerLabelFilters
            },
            all = all
        )
        results = []
        for container in containerList:
            containerLabels = container.attrs.get("Config", {}).get("Labels", {})
            containerType = containerLabels.get("%s.type" % Container.LABEL_PREFIX)
            containerConfig = json.loads(
                containerLabels.get("%s.config" % Container.LABEL_PREFIX),
                object_pairs_hook=collections.OrderedDict
            )
            # app
            if containerType == "application":
                results.append(
                    getApplication(
                        self.getProjectData(),
                        containerConfig
                    )
                )
            # service
            elif containerType == "service":
                results.append(
                    getService(
                        self.getProjectData(),
                        containerConfig
                    )
                )

        # sort results
        def sortKeyName(value):
            return value.getName()
        return sorted(results, key=sortKeyName)

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
        matchedServices = self.dockerFetch("service", name)
        if matchedServices:
            self._services.append(matchedServices[0])
            return matchedServices[0]

        # check path
        if self.path:
            servicesParser = self.getServicesParser()
            if not name:
                serviceNames = servicesParser.getServiceNames()
                if not serviceNames:
                    raise ParserError("No services defined.")
                name = serviceNames[0]
            name = str(name)
            self.logger.debug(
                "Build service handler '%s' for project '%s' from YAML configuration files.",
                name,
                self.getShortUid()
            )
            serviceConfig = servicesParser.getServiceConfiguration(name)
            serviceConfig["_relationships"] = servicesParser.getServiceRelationships(name)
            service = getService(
                self.getProjectData(),
                serviceConfig
            )
            self._services.append(service)
            return service

        # not found
        raise ContainerNotFoundError("Unable to find service with name '%s.'" % name)

    def getApplication(self, name = None):
        """
        Get an application handler.

        :param name: Application name, , if not provided first found application is returned
        :return: Application handler
        :rtype: .application.base.BasePlatformApplication
        """

        # already loaded
        for application in self._applications:
            if application and (not name or application.getName() == name):
                return application   

        # use to sort application list by name
        def appSortKey(value):
            return value.getName()

        # check docker
        matchedApps = self.dockerFetch("application", name)
        if matchedApps:
            self._applications.append(matchedApps[0])
            self._applications = sorted(self._applications, key=appSortKey)
            return matchedApps[0]

        # check path
        if self.path:
            appsParser = self.getApplicationsParser()
            if not name:
                appNames = appsParser.getApplicationNames()
                if not appNames:
                    raise ParserError("No applications defined.")
                name = appNames[0]
            name = str(name)
            self.logger.debug(
                "Build application handler '%s' for project '%s' from YAML configuration files.",
                name,
                self.getShortUid()
            )
            appConfig = appsParser.getApplicationConfiguration(name)
            # init all service dependencies for app
            appServiceDependencies = list(appConfig.get("relationships", {}).values())
            for serivceName in appServiceDependencies:
                self.getService(serivceName.split(":")[0])
            # get app
            app = getApplication(
                self.getProjectData(),
                appConfig
            )
            self._applications.append(app)
            self._applications = sorted(self._applications, key=appSortKey)
            return app

        # not found
        raise ContainerNotFoundError("Unable to find application with name '%s.'" % name)

    def getRouter(self):
        """
        Get main project router.
        """
        if self._router: return self._router
        self._router = PlatformRouter()
        return self._router

    def buildRouterNginxConfig(self, params={}):
        """
        Get Nginx configuration for router.
        """
        # get router
        router = self.getRouter()
        # retrieve all applications
        appParser = self.getApplicationsParser()
        for appName in appParser.getApplicationNames():
            self.getApplication(appName)
        if len(self._applications) == 0:
            raise Exception("Project must contain at least one application.")
        params["_enable_service_routes"] = self.config.get("option_enable_service_routes")
        # generate nginx config
        return router.generateNginxConfig(self._applications, self._services, params=params)        

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
        # generate nginx config
        nginxConfig = self.buildRouterNginxConfig()
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
        # remove router from project network
        dockerClient = PlatformProject.getDockerClient()
        networkName = Container.staticGetNetworkName(self.getProjectData())
        try:
            network = dockerClient.networks.get(networkName)
            network.disconnect(
                router.getContainer()
            )
        except docker.errors.NotFound:
            pass
        except docker.errors.APIError:
            pass
        # restart router
        router.getContainer().restart()

    def purge(self, dryRun = True):
        """
        Purge all docker images and volumes specific to this project.

        :param dryRun: If true do not perform actions, only log affects resources
        """
        # remove all running project containers
        for container in self.dockerFetch(): container.stop()
        # remove from router
        self.removeRouter()
        # get docker client
        dockerClient = self.getDockerClient()
        # remove volumes
        volumeList = dockerClient.volumes.list(
            filters = {
                "label" : [
                    "%s.project-uid=%s" % (Container.LABEL_PREFIX, self.getUid())
                ]
            }
        )
        for volume in volumeList:
            self.logger.info("Delete Docker volume '%s.'" % volume.short_id)
            if not dryRun: volume.remove()
        # remove images
        imageList = dockerClient.images.list(
            filters = {
                "label" : [
                    "%s.project-uid=%s" % (Container.LABEL_PREFIX, self.getUid())
                ]
            }
        )
        for image in imageList:
            self.logger.info("Delete Docker image '%s.'" % image.short_id)
            if not dryRun: dockerClient.images.remove(image.id)
        # remove network
        networkList = dockerClient.networks.list(
            filters = {
                "label" : [
                    "%s.project-uid=%s" % (Container.LABEL_PREFIX, self.getUid())
                ]
            }
        )
        for network in networkList:
            self.logger.info("Delete Docker network '%s.'" % network.short_id)
            if not dryRun: network.remove()

    @staticmethod
    def getAllActiveProjects():
        """ Get list of all currently active projects. """
        dockerClient = PlatformProject.getDockerClient()
        # fetch projects with their network which contains their project uid
        networkList = dockerClient.networks.list(
            filters = {
                "label" : [Container.LABEL_PREFIX]
            }
        )
        projects = []
        for network in networkList:
            labels = network.attrs.get("Labels", {})
            projectUid = labels.get("%s.project-uid" % Container.LABEL_PREFIX)
            if not projectUid: continue
            projects.append(PlatformProject.fromDocker(projectUid))
        return projects
        
    def start(self):
        """ Start project. """
        # start pre-app services
        startGroups = [BasePlatformService.START_PRE_APP_A, BasePlatformService.START_PRE_APP_B]
        servicesParser = self.getServicesParser()
        for startGroup in startGroups:
            for serviceName in servicesParser.getServiceNames():
                service = self.getService(serviceName)
                if service.getStartGroup() == startGroup:
                    service.start()
        # start applications
        applicationsParser = self.getApplicationsParser()
        for applicationName in applicationsParser.getApplicationNames():
            application = self.getApplication(applicationName)
            application.start()
        # start post-app services
        startGroups = [BasePlatformService.START_POST_APP_A, BasePlatformService.START_POST_APP_B]
        for startGroup in startGroups:
            for serviceName in servicesParser.getServiceNames():
                service = self.getService(serviceName)
                if service.getStartGroup() == startGroup:
                    service.start()
        # update router
        self.addRouter()

    def stop(self):
        """ Stop project. """
        projectContainers = self.dockerFetch(all=True)
        for container in projectContainers:
            container.stop()
        self.removeRouter()
