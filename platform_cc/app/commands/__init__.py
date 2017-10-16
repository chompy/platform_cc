import os
import sys
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '../..')))
from app.platform_project import PlatformProject
from app.platform_logger import PlatformLogger

def getLogger(command):
    return PlatformLogger(command)

def getProject(command, withLogger = True):
    projectPath = command.option("path")
    return PlatformProject(
        projectPath if projectPath else "",
        getLogger(command) if withLogger else None
    )

def getAppsToInvoke(command):
    appInvokeList = command.option("apps")
    project = getProject(command)
    apps = project.getApplications()
    filteredApps = []
    for app in apps:
        if not appInvokeList or app.config.getName().lower() in appInvokeList:
            filteredApps.append(app)
    return filteredApps