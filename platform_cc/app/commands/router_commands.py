from cleo import Command
from app.commands import getProject, getAppsToInvoke, getLogger
from app.platform_router import PlatformRouter

class RouterStart(Command):
    """
    Start router.

    router:start
    """

    def handle(self):
        PlatformRouter(getLogger(self)).start()


class RouterStop(Command):
    """
    Stop router.

    router:stop
    """

    def handle(self):
        PlatformRouter(getLogger(self)).stop()

class RouterAdd(Command):
    """
    Add project to router.

    router:add
    {--p|path=? : Path to project root. (Default=current directory)}
    """
    
    def handle(self):
        project = getProject(self)
        PlatformRouter(getLogger(self)).addProject(
            project
        )

class RouterRemove(Command):
    """
    Remove project from router.

    router:remove
    {--p|path=? : Path to project root. (Default=current directory)}
    """
    
    def handle(self):
        project = getProject(self)
        PlatformRouter(getLogger(self)).removeProject(
            project
        )