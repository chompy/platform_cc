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

    def run(self):
        """
        Perform task.
        """
        pass