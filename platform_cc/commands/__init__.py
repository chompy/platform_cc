import os
from project import PlatformProject

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
    