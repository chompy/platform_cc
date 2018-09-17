import logging

class BaseTaskHandler:
    """
    Base task handler class.
    """

    def __init__(self, project, params = {}):
        self.project = project
        self.params = dict(params)
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

        return (
            self.project.getApplication(appName),
            appPath
        )
            

    def run(self):
        """
        Perform task.
        """
        pass