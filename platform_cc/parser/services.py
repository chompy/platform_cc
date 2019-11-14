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
import yaml
import yamlordereddictloader
import collections
from .base import BasePlatformParser
from .yaml_archive import ArchiveTag
from .yaml_include import IncludeTag
from platform_cc.exception.parser_error import ParserError


class ServicesParser(BasePlatformParser):
    """
    Services (.platform/services.yaml) parser.
    """

    """ Paths to services yaml file. """
    YAML_PATHS = [
        ".platform/services.yaml",
        ".platform/services.pcc.yaml"
    ]

    def __init__(self, projectPath):
        BasePlatformParser.__init__(self, projectPath)
        yamlPaths = []
        for yamlPath in self.YAML_PATHS:
            yamlFullPath = os.path.join(
                self.projectPath,
                yamlPath
            )
            if os.path.isfile(yamlFullPath):
                yamlPaths.append(yamlFullPath)
        self.services = self._readYamls(yamlPaths)

    def _readYaml(self, path):
        loadConf = {}
        ArchiveTag.base_path = os.path.dirname(path)
        IncludeTag.base_path = ArchiveTag.base_path
        with open(path, "r") as f:
            loadConf = yaml.load(f, Loader=yamlordereddictloader.SafeLoader)
            for key in loadConf:
                if type(loadConf[key]) is not collections.OrderedDict:
                    continue
                if "_from" not in loadConf[key]:
                    loadConf[key]["_from"] = []
                if path not in loadConf[key]["_from"]:
                    loadConf[key]["_from"].append(path)
        return loadConf

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
        if name not in self.services:
            raise ParserError(
                "Service '%s' is not defined." % name
            )
        # ensure 'type' parameter is present
        if "type" not in self.services[name]:
            raise ParserError(
                """
                'type' configuration parameter is missing for service '%s.'
                """ % (
                    name
                )
            )

        return self.services[name].get("type", "")

    def getServiceRelationships(self, name):
        """
        Get relationships for given service.

        :param name: Service name
        :return: Dictionary with relationships
        :rtype: dict
        """
        name = str(name)
        if name not in self.services:
            raise ParserError("Service '%s' is not defined." % name)
        return self.services[name].get("relationships", {})

    def getServiceConfiguration(self, name):
        """
        Get dictionary containing configuration for a
        given service.

        :param name: Service name
        :return: Dictionary with service configuration
        :rtype: dict
        """
        name = str(name)
        if name not in self.services:
            raise ParserError("Service '%s' is not defined." % name)
        serviceConf = self.services[name].get("configuration", {}).copy()
        serviceConf["_name"] = name
        serviceConf["_type"] = self.getServiceType(name)
        serviceConf["_disk"] = self.services[name].get("disk", 0)
        serviceConf["_is_default_config"] = os.path.basename(
            self.services[name].get("_path", "")
        ) == os.path.basename(self.YAML_PATHS[0])
        return serviceConf
