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

from __future__ import absolute_import
import os
import json
from platform_cc.project import PlatformProject
from terminaltables import SingleTable

def getProject(command):
    """
    Get PlatformProject object for current command.

    :param command: Current command being ran
    :return: Project
    :rtype: PlatformProject
    """
    path = command.option("path")
    if not path: path = os.getcwd()
    return PlatformProject(path)
    
def outputTable(command, title, data):
    """
    Output an ASCII table to the terminal.

    :param command: Current command being ran
    :param title: Table title
    :param data: Table data
    """
    table = SingleTable(
        data,
        title
    )
    command.line(table.table)

def outputJson(command, data):
    """
    Output JSON to the terminal.

    :param command: Current command being ran
    :param data: Data to display as JSON
    """
    command.line(
        json.dumps(
            data,
            indent = 4
        )
    )