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
from cleo import Command
from platform_cc.commands import getProject, outputJson, outputTable

class ServiceRestart(Command):
    """
    Restart one or more services.

    service:restart
        {name* : Name(s) of service.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        for name in self.argument("name"):
            service = project.getService(name)
            service.restart()

class ServiceList(Command):
    """
    List all services.

    service:list
        {--p|path=? : Path to project root. (Default=current directory)}
        {--u|uid=? : Project uid.}
        {--j|json : If set output in JSON.}
    """

    def handle(self):
        project = getProject(self)
        services = project.dockerFetch("service", all=True)
        # build service data
        serviceData = {}
        for service in services:
            serviceContainer = service.getContainer()
            serviceData[service.getName()] = {
                "id"                  : serviceContainer.id if serviceContainer else "",
                "name"                : service.getName(),
                "type"                : service.getType(),
                "container_name"      : service.getContainerName(),
                "base_docker_image"   : service.getBaseImage(),
                "docker_image"        : service.getDockerImage(),
                "ip_address"          : service.getContainerIpAddress()
            }

        # json output
        if self.option("json"):
            outputJson(
                self,
                serviceData
            )
            return
        
        # terminal tables output
        tableData = [
            ("Name", "Type", "Container", "Image", "IP"),
        ]
        for service in serviceData:
            tableData.append(
                (
                    serviceData[service]["name"],
                    serviceData[service]["type"],
                    serviceData[service]["container_name"],
                    serviceData[service]["docker_image"],
                    serviceData[service]["ip_address"] if serviceData[service]["ip_address"] else "N/A"
                )
            )
        outputTable(
            self,
            "Project '%s' - Services" % project.getUid()[0:6],
            tableData
        )

class ServiceShell(Command):
    """
    Shell in to a service.

    service:shell
        {name : Name of service.}
        {--p|path=? : Path to project root. (Default=current directory)}
        {--u|uid=? : Project uid.}
        {--c|command=? : Command to run. (Default=bash)}
        {--user=? : User to run command as. (Default=root)}
    """

    def handle(self):
        project = getProject(self)
        service = project.getService(
            self.argument("name")
        )
        command = self.option("command")
        if not command: command = "bash"
        user = self.option("user")
        if not user: user = "root"
        service.shell(
            command,
            user
        )        