class BasePlatformVariables:
    """
    Base class for handling storage of variables.
    """

    def __init__(self, projectPath, projectConfig = {}):
        """
        Constructor.

        :param projectPath: Path to project root
        :param projectConfig: Project configuration
        """
        self.projectPath = str(projectPath)
        self.projectConfig = dict(projectConfig)

    def set(self, key, value):
        """
        Set a project variable.

        :param key: Name of variable
        :param value: Value of variable
        """
        pass

    def get(self, key, default = None):
        """
        Retrieve a project variable.

        :param key: Name of variable
        :param default: Default value to return if variable does not exist
        :return: Retrieved variable value
        """
        pass

    def delete(self, key):
        """
        Delete a project variable.

        :param key: Name of variable
        """
        pass

    def all(self):
        """
        Retrieve dictionary containing all project variables.

        :return: Dictionary containing all variables
        :rtype: dict
        """
        pass