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
        applicationsParser = project.getApplicationsParser()
        for applicationName in applicationsParser.getApplicationNames():
            application = project.getApplication(applicationName)
            application.start()
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
        servicesParser = project.getServicesParser()
        for serviceName in servicesParser.getServiceNames():
            service = project.getService(serviceName)
            service.stop()
        project.removeRouter()

class ProjectRestart(Command):
    """
    Restart all applications and services in a project.

    project:restart
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        applicationsParser = project.getApplicationsParser()
        for applicationName in applicationsParser.getApplicationNames():
            application = project.getApplication(applicationName)
            application.restart()
        servicesParser = project.getServicesParser()
        for serviceName in servicesParser.getServiceNames():
            service = project.getService(serviceName)
            service.restart()
        project.addRouter()

class ProjectRoutes(Command):
    """
    List routes in project.

    project:routes
        {--p|path=? : Path to project root. (Default=current directory)}
        {--j|json : If set output in JSON.}
    """

    def handle(self):
        project = getProject(self)
        routesParser = project.getRoutesParser()
        routes = routesParser.getRoutes()

        # json output
        if self.option("json"):
            outputJson(
                self,
                routes
            )
            return

        # terminal tables output
        tableData = [
            ("Route", "Type", "Upstream/To"),
        ]
        for route in routes:
            for host in route.get("hostnames", []):
                routeName = "%s://%s/%s" % (
                    route.get("scheme", "http"),
                    host,
                    route.get("path", "").lstrip("/")
                )
                tableData.append(
                    (
                        routeName,
                        route.get("type", "upstream"),
                        route.get("upstream", "") if route.get("type") == "upstream" else route.get("to", "")
                    )
                )
        outputTable(
            self,
            "Project '%s' - Routes" % project.getUid()[0:6],
            tableData
        )