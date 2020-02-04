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
from .common import outputJson, outputTable
from ..exception.state_error import StateError
from ..core.project import PlatformProject
from ..core.router import PlatformRouter
from ..application.base import BasePlatformApplication
from ..services.base import BasePlatformService

class AllStop(Command):
    """
    Stop all running projects.

    all:stop
    """

    def handle(self):
        projects = PlatformProject.getAllActiveProjects()
        for project in projects:
            projectContainers = project.dockerFetch(all=True)
            for container in projectContainers:
                container.stop()
            project.removeRouter()
        PlatformRouter().stop()

class AllPurge(Command):
    """
    Purge all Platform.CC related Docker contains, images, volumes, and networks.

    all:purge
        {--d|dry-run : List items to be purged.}
    """

    def handle(self):
        dryRun = bool(self.option("dry-run"))
        # inform user and wait 5 seconds
        if not dryRun:
            self.line(
                "<question>!!! Purge of ALL Platform.CC projects will commence in 10 seconds. Press CTRL+C to cancel. !!!</question>"
            )
            time.sleep(10)
        projects = PlatformProject.getAllActiveProjects()
        for project in projects:
            project.purge(dryRun)

class AllList(Command):
    """
    List details about all running projects.

    all:list
        {--j|json : If set output in JSON.}
    """

    def handle(self):
        projects = PlatformProject.getAllActiveProjects()
        outputData = []
        for project in projects:
            appList = []
            serviceList = []
            subnet = ""
            projectContainers = project.dockerFetch(all=True)
            for container in projectContainers:
                if isinstance(container, BasePlatformApplication):
                    appList.append({
                        "name" : container.getName(),
                        "type" : container.getType(),
                        "image" : container.getDockerImage()
                    })
                elif isinstance(container, BasePlatformService):
                    serviceList.append({
                        "name" : container.getName(),
                        "type" : container.getType(),
                        "image" : container.getDockerImage()
                    })
                if not subnet:
                    subnet = container.getNetwork().attrs.get("IPAM", {}).get("Config", [{}])[0].get("Subnet", "n/a")
            if not appList and not serviceList: continue
            outputData.append({
                "short_uid" : project.getShortUid(),
                "applications" : appList,
                "services" : serviceList,
                "subnet" : subnet
            })

        # json output
        if self.option("json"):
            outputJson(
                self,
                outputData
            )
            return
        
        # terminal tables output
        tableData = [
            ("Uid", "Applications", "Services", "Subnet"),
        ]
        for projectData in outputData:
            appNames = []
            serviceNames = []
            for app in projectData.get("applications", []):
                appNames.append("%s/%s" % (app["name"], app["type"]))
            for service in projectData.get("services", []):
                serviceNames.append("%s/%s" % (service["name"], service["type"]))
            tableData.append(
                (
                    projectData["short_uid"],
                    "\n".join(appNames).strip(),
                    "\n".join(serviceNames).strip(),
                    projectData["subnet"]
                )
            )
        outputTable(
            self,
            "Active Projects",
            tableData
        )

