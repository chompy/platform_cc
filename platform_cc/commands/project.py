"""
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
"""

from __future__ import absolute_import
import os
import time
import docker
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
        {--u|uid=? : Project uid.}
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

class ProjectOptionSet(Command):
    """
    Enable or disable project option.

    project:option_set
        {key : Name of option to set.}
        {enable : Whether or not to enable or disable option.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        enabled = self.argument("enable").lower() in ["1", "yes", "y", "enabled", "enable", "true", "t"]
        project.config.set(
            "option_%s" % self.argument("key").lower(),
            "enabled" if enabled else ""
        )

class ProjectOptionList(Command):
    """
    List project options.

    project:options
        {--p|path=? : Path to project root. (Default=current directory)}
        {--j|json : If set output in JSON.}
    """

    def handle(self):
        project = getProject(self)
        options = [
            {
                "name"          : "use_mount_volumes",
                "description"   : "Use Docker volumes for application mount volumes.",
                "enabled"       : bool(project.config.get("option_use_mount_volumes", False))
            },
            {
                "name"          : "enable_cron",
                "description"   : "Enable Cron jobs inside application containers.",
                "enabled"       : bool(project.config.get("option_enable_cron", False))
            }
        ]
        # json output
        if self.option("json"):
            outputJson(self, options)
            return

        # terminal tables output
        tableData = [
            ("Name", "Description", "Enabled")
        ]
        for option in options:
            tableData.append(
                (
                    option.get("name").upper(),
                    option.get("description"),
                    option.get("enabled")
                )
            )
        outputTable(
            self,
            "Project '%s' - Options" % project.getUid()[0:6],
            tableData
        )

class ProjectPurge(Command):
    """
    Purge all docker images and volumes specific to this project.

    project:purge
        {--d|dry-run : List items to be purged.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        dryRun = bool(self.option("dry-run"))
        # inform user and wait 5 seconds
        if not dryRun:
            self.line(
                "<question>!!! Purge of project '%s' will commence in 5 seconds. Press CTRL+C to cancel. !!!</question>" % (
                    project.getShortUid()
                )
            )
            time.sleep(5)
        # services
        serviceParser = project.getServicesParser()
        for serviceName in serviceParser.getServiceNames():
            service = project.getService(serviceName)
            service.purge(dryRun)
        # applications
        appParser = project.getApplicationsParser()
        for appName in appParser.getApplicationNames():
            app = project.getApplication(appName)
            app.purge(dryRun)
        # remove from router
        if not dryRun:
            project.removeRouter()
        # delete network
        app = project.getApplication(appParser.getApplicationNames()[0])
        networkName = app.getNetworkName()
        try:
            network = app.docker.networks.get(networkName)
            if not dryRun:
                try:
                    network.disconnect(project.getRouter().getContainerName())
                except docker.errors.APIError:
                    pass
                network.remove()
            app.logger.info(
                "Deleted '%s' Docker network.",
                networkName
            )
        except docker.errors.NotFound:
            pass
        