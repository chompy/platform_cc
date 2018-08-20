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
from .base import BasePlatformParser
from platform_cc.exception.parser_error import ParserError

class ApplicationsParser(BasePlatformParser):
    """
    Applications (.platform.app.yaml) parser.
    """

    """ Filenames of application configuration yaml file. """
    YAML_FILENAMES = [
        ".platform.app.yaml",
        ".platform.app.pcc.yaml"
    ]

    def __init__(self, projectPath):
        BasePlatformParser.__init__(self, projectPath)
        self._compile()

    def _compile(self):
        """
        Compile all application yamls into dictionary.
        """
        self.applications = {}
        yamlPaths = self.getYamlPaths()
        if not yamlPaths or len(yamlPaths) < 1:
            raise ParserError("No applications have been defined.")
        for yamlPath in yamlPaths:
            appConfig = self._readYamls(yamlPath)
            name = appConfig.get(
                "name",
                os.path.basename(os.path.dirname(yamlPath[0]))
            )
            # skip if no name or name already exists
            if not name or name in self.applications:
                continue
            self.applications[name] = appConfig

    def getYamlPaths(self):
        """
        Get path to all application yaml files.

        :return: List of application yaml files
        :rtype: list
        """
        yamlList = []
        # get all application yaml files in a directory
        def getYamlPathsInDir(path):
            if not os.path.isdir(path): return []
            yamlFiles = []
            for yamlFilename in self.YAML_FILENAMES:
                yamlPath = os.path.join(path, yamlFilename)
                if os.path.isfile(yamlPath):
                    yamlFiles.append(yamlPath)
            return yamlFiles
        # find app yaml in project path
        yamlDirList = getYamlPathsInDir(self.projectPath)
        if yamlDirList:
            yamlList.append(yamlDirList)
        # find app yaml in sub directories
        for name in os.listdir(self.projectPath):
            fullPath = os.path.join(
                self.projectPath,
                name
            )
            if os.path.isdir(fullPath):
                yamlDirList = getYamlPathsInDir(fullPath)
                if yamlDirList:
                    yamlList.append(yamlDirList)
        return yamlList

    def getApplicationNames(self):
        """
        Get list of names of all applications.

        :return: List of application names
        :rtype: list
        """
        return list(self.applications.keys())

    def getApplicationConfiguration(self, name):
        """
        Get dictionary containing configuration for a
        given application.

        :param name: Application name
        :return: Dictionary with application configuration
        :rtype: dict
        """
        name = str(name)
        if not name in self.applications:
            raise ParserError("Application '%s' is not defined." % name)
        applicationConfig = self.applications[name].copy()
        return applicationConfig