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