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

class BasePlatformVariables:
    """
    Base class for handling storage of variables.
    """

    def __init__(self, projectPath, projectConfig = None):
        """
        Constructor.

        :param projectPath: Path to project root
        :param projectConfig: Project configuration
        """
        self.projectPath = str(projectPath)
        self.projectConfig = projectConfig

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