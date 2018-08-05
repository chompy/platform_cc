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

from .var_json_file import JsonFileVariables
from .var_dict import DictVariables

def getVariableStorage(config = {}):
    """
    Init a variable storage class from project configuration.

    :param config: Dict containing configuration
    """

    # enforce dict
    config = dict(config)

    # determine which storage handler to use
    variableStorageHandler = str(
        config.get("storage_handler", "json_file")
    )
    
    # use json file handler
    if variableStorageHandler in ["json_file", "json"]:
        return JsonFileVariables(config)
    elif variableStorageHandler in ["dict", "dictionary"]:
        return DictVariables(config)

    # not found
    raise NotImplementedError(
        "Variable storage handler '%s' has not been implemented." % (
            variableStorageHandler
        )
    )
    