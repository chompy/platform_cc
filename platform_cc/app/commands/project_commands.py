from cleo import Command
from app.commands import getProject, getAppsToInvoke, getLogger
from app.platform_router import PlatformRouter

class ProjectStart(Command):
    """
    Start one or more applications in a project.

    project:start
        {--a|apps=* : Comma delimited list of applications to start. (Default=all)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        for app in getAppsToInvoke(self):
            app.start()
        router = PlatformRouter(getLogger(self))
        router.start()
        router.addProject(getProject(self))

class ProjectStop(Command):
    """
    Stop one or more applications in a project.

    project:stop
        {--a|apps=* : Comma delimited list of applications to stop. (Default=all)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        router = PlatformRouter(getLogger(self))
        router.removeProject(getProject(self))
        for app in getAppsToInvoke(self):
            app.stop()

class ProjectBuild(Command):
    """
    Perform pre-build actions and run build hooks.

    project:build
        {--a|apps=* : Comma delimited list of applications to build. (Default=all)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        for app in getAppsToInvoke(self):
            app.build()

class ProjectDeploy(Command):
    """
    Run deploy hooks.

    project:deploy
        {--a|apps=* : Comma delimited list of applications to deploy. (Default=all)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        for app in getAppsToInvoke(self):
            app.deploy()

class ProjectInfo(Command):

    """
    Display information about the project.

    project:info
        {--p|path=? : Path to project root. (Default=current directory)}
    """
    def handle(self):
        project = getProject(
            self,
            True
        )
        project.outputInfo()