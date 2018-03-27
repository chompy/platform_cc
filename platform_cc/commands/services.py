import os
from cleo import Command
from commands import getProject

def getService(command):
    """
    Get service from command.
    """
    project = getProject(command)
    serviceName = command.argument("name")
    project.servicesParser

    return getService()

class ServiceStart(Command):
    """
    Start a service.

    service:start
        {name : Name of service.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        service = project.getService(
            self.argument("name")
        )
        service.start()

class ServiceStop(Command):
    """
    Stop a service.

    service:stop
        {name : Name of service.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        service = project.getService(
            self.argument("name")
        )
        service.stop()

class ServiceRestart(Command):
    """
    Restart a service.

    service:restart
        {name : Name of service.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        service = project.getService(
            self.argument("name")
        )
        service.restart()