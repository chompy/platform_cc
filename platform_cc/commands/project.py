from __future__ import absolute_import
import os
from cleo import Command
from platform_cc.commands import getProject, outputJson, outputTable

class ProjectStart(Command):
    """
    Start all applications and services in a project.

    project:start
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        servicesParser = project.getServicesParser()
        for serviceName in servicesParser.getServiceNames():
            service = project.getService(serviceName)
            service.start()
            self.line(service.getContainerName())
        applicationsParser = project.getApplicationsParser()
        for applicationName in applicationsParser.getApplicationNames():
            application = project.getApplication(applicationName)
            application.start()
            self.line(application.getContainerName())
        project.addRouter()

class ProjectStop(Command):
    """
    Stop all applications and services in a project.

    project:stop
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        applicationsParser = project.getApplicationsParser()
        for applicationName in applicationsParser.getApplicationNames():
            application = project.getApplication(applicationName)
            application.stop()
            self.line(application.getContainerName())
        servicesParser = project.getServicesParser()
        for serviceName in servicesParser.getServiceNames():
            service = project.getService(serviceName)
            service.stop()
            self.line(service.getContainerName())
        project.removeRouter()