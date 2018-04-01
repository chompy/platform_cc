import os
import json
import time
import hashlib
import base36
import random
from variables import getVariableStorage
from variables.var_json import JsonVariables
from parser.services import ServicesParser
from parser.applications import ApplicationsParser
from services import getService
from application import getApplication

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

    def _getProjectData(self):
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
            "path"      : self.path,
            "uid"       : self.getUid(),
            "short_uid" : self.getUid()[0:6],
            "entropy"   : self.getEntropy(),
            "config"    : self.config.all(),
            "variables" : self.variables.all(),
            "services"  : serviceData
        }

    def getServicesParser(self):
        """
        Get services parser.

        :return: Services parser
        :rtype: .parser.services.ServicesParser
        """
        if self._servicesParser:
            return self._servicesParser
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
        self._applicationsParser = ApplicationsParser(self.path)
        return self._applicationsParser

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
        service = getService(
            self._getProjectData(),
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
        application = getApplication(
            self._getProjectData(),
            appConfig
        )
        self._applications.append(application)
        return application