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

import os
import time
from cleo import Command, Output
from platform_cc.commands import getProject, outputJson, outputTable
from platform_cc.exception.state_error import StateError

class ApplicationRestart(Command):
    """
    Restart one or more applications.

    application:restart
        {name* : Name(s) of application.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        for name in self.argument("name"):
            application = project.getApplication(name)
            application.restart()

class ApplicationList(Command):
    """
    List all applications.

    application:list
        {--p|path=? : Path to project root. (Default=current directory)}
        {--j|json : If set output in JSON.}
    """

    def handle(self):
        project = getProject(self)
        appNames = project.getApplicationsParser().getApplicationNames()

        # build app data
        appData = {}
        for appName in appNames:
            application = project.getApplication(
                appName
            )
            appContainer = application.getContainer()
            appData[appName] = {
                "id"                  : appContainer.id if appContainer else "",
                "name"                : application.getName(),
                "container_name"      : application.getContainerName(),
                "type"                : application.getType(),
                "base_docker_image"   : application.getBaseImage(),
                "docker_image"        : application.getDockerImage(),
                "ip_address"          : application.getContainerIpAddress()
            }

        # json output
        if self.option("json"):
            outputJson(
                self,
                appData
            )
            return
        
        # terminal tables output
        tableData = [
            ("Name", "Type", "Container", "Image", "IP"),
        ]
        for application in appData:
            tableData.append(
                (
                    appData[application]["name"],
                    appData[application]["type"],
                    appData[application]["container_name"],
                    appData[application]["docker_image"],
                    appData[application]["ip_address"] if appData[application]["ip_address"] else "N/A"
                )
            )
        outputTable(
            self,
            "Project '%s' - Applications" % project.getUid()[0:6],
            tableData
        )

class ApplicationShell(Command):
    """
    Shell in to an application.

    application:shell
        {--name=? : Name of application. (Default=first available application)}
        {--p|path=? : Path to project root. (Default=current directory)}
        {--u|uid=? : Project uid.}
        {--c|command=? : Command to run. (Default=bash)}
        {--user=? : User to run command as. (Default=web)}
    """

    def handle(self):
        project = getProject(self)
        app = project.getApplication(
            self.option("name")
        )
        command = self.option("command")
        if not command: command = "bash"
        user = self.option("user")
        if not user: user = "web"
        app.shell(
            command,
            user
        )        

class ApplicationBuild(Command):
    """
    (Re)Build application.

    application:build
        {--name=? : Name of application. (Default=first available application)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        name = self.option("name")
        if not name:
            name = project.getApplicationsParser().getApplicationNames()[0]
        application = project.getApplication(name)
        # inform user and wait 5 seconds
        self.line(
            "<question>!!! Rebuild of application '%s' will commence in 5 seconds. Press CTRL+C to cancel. !!!</question>" % (
                application.getName()
            )
        )
        time.sleep(5)
        output = application.build()
        if self.output.get_verbosity() >= Output.VERBOSITY_VERBOSE:
            self.line(output)

class ApplicationDeployHook(Command):
    """
    Run deploy hook for application.

    application:deploy
        {--name=? : Name of application. (Default=first available application)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        name = self.option("name")
        if not name:
            name = project.getApplicationsParser().getApplicationNames()[0]
        application = project.getApplication(name)
        output = application.deploy()
        if self.output.get_verbosity() >= Output.VERBOSITY_VERBOSE:
            self.line(output)

class ApplicationPull(Command):
    """
    Pull base application image.

    application:pull
        {--name=? : Name of application. (Default=first available application)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        name = self.option("name")
        if not name:
            name = project.getApplicationsParser().getApplicationNames()[0]
        application = project.getApplication(name)
        application.pullImage()