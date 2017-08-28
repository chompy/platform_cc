from app.platform_project import PlatformProject
from app.platform_logger import PlatformLogger

def getLogger(command):
    return PlatformLogger(command)

def getProject(command):
    projectPath = command.option("path")
    return PlatformProject(
        projectPath if projectPath else "",
        getLogger(command)
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