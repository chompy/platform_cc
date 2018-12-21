import os
import logging
from ...exception.container_command_error import ContainerCommandError

class BaseTaskHandler:
    """
    Base task handler class.
    """

    def __init__(self, project, params = {}):
        self.project = project
        self.params = dict(params)
        self.resolveShortCodes()
        self.logger = logging.getLogger(__name__)

    @classmethod
    def getType(cls):
        """
        Get type name of task handler.
    
        :rtype: string
        """
        return ""

    def checkParams(self, names = []):
        """
        Check a parameter value key exists for each name
        provided. Raise ValueError if not.

        :param list: List of dictionary keys
        """
        for name in names:
            if not self.params.has_key(name):
                raise ValueError("Task handler '%s' is missing '%s' parameter." % (self.getType(), name))            

    def _replaceShortCodes(self, string):
        """
        Replace short codes in given string with appropiate
        values.

        :param string: String with short codes
        :rtype: string
        """
        # project_dirname shortcode
        # get the name of the directory the project
        if "{PROJECT_DIRNAME}" in string:
            string = string.replace(
                "{PROJECT_DIRNAME}",
                os.path.basename(self.project.path)
            )
        # project_path
        # get path to project
        if "{PROJECT_PATH}" in string:
            string = string.replace(
                "{PROJECT_PATH}",
                self.project.path
            )
        # project_application
        # get name of first available application
        if "{PROJECT_APPLICATION}" in string:
            string = string.replace(
                "{PROJECT_APPLICATION}",
                self.project.getApplication().getName()
            )
        return string

    def resolveShortCodes(self):
        """
        Replace all shortcodes found in task parameters.
        """
        for key in self.params: 
            if type(self.params[key]) is str:
                self.params[key] = self._replaceShortCodes(self.params[key])
            elif type(self.params[key]) is list:
                for index in range(len(self.params[key])):
                    self.params[key][index] = self._replaceShortCodes(self.params[key][index])

    def checkCondition(self):
        """
        Evaluate condition parameter.

        :rtype: bool
        """
        conditionParams = self.params.get("condition")
        if not conditionParams: return True
        appName = ""
        conditionCommand = ""
        if type(conditionParams) is str:
            conditionCommand = conditionParams
        elif type(conditionParams) is list:
            if len(conditionParams) == 1:
                conditionCommand = conditionParams[0]
            elif len(conditionParams) >= 2:
                appName = conditionParams[0]
                conditionCommand = conditionParams[1]
        elif type(conditionParams) is dict:
            appName = conditionParams.get("application", "")
            conditionCommand = conditionParams.get("command", "")
        
        application = self.project.getApplication(appName)
        try:
            application.runCommand(conditionCommand)
        except ContainerCommandError:
            return False
        return True

    def parseAppPath(self, path):
        """
        Parse a path string that points to a path inside an
        application container.
        Expected syntax... "app_name:path" or "path"

        :param string: Path to
        :rtype: tuple
        :return: Tuple containing application and path
        """
        pathSplit = path.strip().split(":")
        appName = None
        appPath = ""
        if len(pathSplit) == 1:
            appPath = pathSplit[0].strip()
        elif len(pathSplit) > 1:
            appName = pathSplit[0].strip()
            appPath = pathSplit[1].strip()

        if "{PROJECT_APPLICATION}" in appName:
            appName = appName.replace(
                "{PROJECT_APPLICATION}",
                self.project.getApplication().getName()
            )

        return (
            self.project.getApplication(appName),
            appPath
        )
            

    def run(self):
        """
        Perform task.
        """
        pass