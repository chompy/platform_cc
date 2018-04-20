import os
from .base import BasePlatformParser
from platform_cc.exception.parser_error import ParserError

class ServicesParser(BasePlatformParser):
    """
    Services (.platform/services.yaml) parser.
    """

    """ Path to services yaml file. """
    YAML_PATH = ".platform/services.yaml"

    def __init__(self, projectPath):
        BasePlatformParser.__init__(self, projectPath)
        yamlPath = os.path.join(
            self.projectPath,
            self.YAML_PATH
        )
        self.services = self._readYaml(yamlPath)

    def getServiceNames(self):
        """
        Get list of all service names for this project.

        :return: List of service names
        :rtype: list
        """
        return list(self.services.keys())

    def getServiceType(self, name):
        """
        Get service type given name of service.

        :param name: Service name
        :return: Service type, name of service software
        :rtype: str
        """
        name = str(name)
        if not name in self.services:
            raise ParserError(
                "Service '%s' is not defined." % name
            )
        # ensure 'type' parameter is present
        if not "type" in self.services[name]:
            raise ParserError(
                "'type' configuration parameter is missing for service '%s.'" % (
                    name
                )
            )

        return self.services[name].get("type", "")

    def getServiceConfiguration(self, name):
        """
        Get dictionary containing configuration for a
        given service.

        :param name: Service name
        :return: Dictionary with service configuration
        :rtype: dict
        """
        name = str(name)
        if not name in self.services:
            raise ParserError("Service '%s' is not defined." % name)
        serviceConf = self.services[name].get("configuration", {}).copy()
        serviceConf["_name"] = name
        serviceConf["_type"] = self.getServiceType(name)
        return serviceConf

