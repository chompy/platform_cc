import os
from cleo import Command
from router import PlatformRouter
from commands import getProject

class RouterStart(Command):
    """
    Start the router.

    router:start
    """

    def handle(self):
        router = PlatformRouter()
        router.start()
        self.line(router.getContainerName())

class RouterStop(Command):
    """
    Stop the router.

    router:stop
    """

    def handle(self):
        router = PlatformRouter()
        router.stop()
        self.line(router.getContainerName())

class RouterAdd(Command):
    """
    Add project to router.

    router:add
    {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        router = PlatformRouter()
        router.addProject(project.getProjectData())
