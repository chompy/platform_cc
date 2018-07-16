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

from .php import PhpApplication

""" Map application names to their application class. """
APPLICATION_MAP = {
    "php"           : PhpApplication
}

def getApplication(project, config):
    """
    Get application handler from configuration.
    
    :param project: Dictionary containing project data
    :param config: Dictionary containing application configuration
    :return: Application handler object
    :rtype: .base.BasePlatformApplication
    """

    # validate config
    if not isinstance(config, dict):
        raise ValueError("Config parameter must be a dictionary (dict) object.")
    if "type" not in config or "name" not in config:
        raise ValueError("Config parameter is missing parameters required for an application.")
    appType = config["type"].split(":")[0]
    if appType not in APPLICATION_MAP:
        raise NotImplementedError(
            "No appliocation handler available for '%s.'" % (
                config["type"]
            )
        )
    # init application
    return APPLICATION_MAP[appType](
        project,
        config
    )