from cleo import Command
from app.commands import getProject, getLogger
from app.platform_router import PlatformRouter

class ProjectStart(Command):
    """
    Start one or more applications in a project.

    project:start
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        project.start()
        router = PlatformRouter(getLogger(self))
        router.start()
        router.addProject(project)

class ProjectStop(Command):
    """
    Stop one or more applications in a project.

    project:stop
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        router = PlatformRouter(getLogger(self))
        router.removeProject(project)
        project.stop()

class ProjectProvision(Command):
    """
    (Re)provisions all apps in a project.

    project:provision
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        getProject(self).provision()

class ProjectDeploy(Command):
    """
    Run deploy hooks.

    project:deploy
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        getProject(self).deploy()

class ProjectInfo(Command):

    """
    Display information about the project.

    project:info
        {--p|path=? : Path to project root. (Default=current directory)}
    """
    def handle(self):
        getProject(self).outputInfo()

class ProjectPurge(Command):

    """
    Purge project. Stop all containers and delete all volumes.

    project:purge
        {--p|path=? : Path to project root. (Default=current directory)}
    """
    def handle(self):
        getProject(self).purge()
