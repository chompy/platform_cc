from __future__ import absolute_import
import os
from cleo import Command
from platform_cc.router import PlatformRouter
from platform_cc.commands import getProject

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
        project.addRouter()

class RouterRemove(Command):
    """
    Remove project from router.

    router:remove
    {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        project.removeRouter()