import os
from .base import BasePlatformParser
from exception.parser_error import ParserError

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

    def getServiceType(self, name):
        """
        Get service type given name of service.

        :param name: Service name
        :return: Service type, name of service software
        :rtype: str
        """
        if not name in self.services:
            raise ValueError("Service '%s' is not defined." % name)
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
        if not name in self.services:
            raise ValueError("Service '%s' is not defined." % name)
        return self.services[name].get("configuration", {})
