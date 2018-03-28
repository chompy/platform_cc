import os
from .base import BasePlatformParser

class ApplicationsParser(BasePlatformParser):
    """
    Applications (.platform.app.yaml) parser.
    """

    """ Filename of application configuration yaml file. """
    YAML_FILENAME = ".platform.app.yaml"

    def __init__(self, projectPath):
        BasePlatformParser.__init__(self, projectPath)
        self._compile()

    def _compile(self):
        """
        Compile all application yamls into dictionary.
        """
        self.applications = {}
        yamlPaths = self.getYamlPaths()
        for yamlPath in yamlPaths:
            appConfig = self._readYaml(yamlPath)
            name = appConfig.get(
                "name",
                os.path.basename(os.path.dirname(yamlPath))
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
        for name in os.listdir(self.projectPath):
            fullPath = os.path.join(
                self.projectPath,
                name
            )
            if os.path.isdir(fullPath):
                yamlPath = os.path.join(fullPath, self.YAML_FILENAME)
                if os.path.isfile(yamlPath):
                    yamlList.append(yamlPath)
            elif name == self.YAML_FILENAME:
                yamlList.append(fullPath)
        return yamlList