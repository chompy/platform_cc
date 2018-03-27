import os
from cleo import Command
from commands import getProject, outputJson, outputTable

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
        self.line(service.getContainerName())

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
        self.line(service.getContainerName())

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
        self.line(service.getContainerName())

class ServiceList(Command):
    """
    List all services.

    service:list
        {--p|path=? : Path to project root. (Default=current directory)}
        {--j|json : If set output in JSON.}
    """

    def handle(self):
        project = getProject(self)
        serviceNames = project.servicesParser.getServiceNames()

        # build service data
        serviceData = {}
        for serviceName in serviceNames:
            service = project.getService(
                serviceName
            )
            serviceContainer = service.getContainer()
            serviceData[serviceName] = {
                "id"                  : serviceContainer.id if serviceContainer else "",
                "name"                : service.getName(),
                "container_name"      : service.getContainerName(),
                "docker_image"        : service.getDockerImage(),
                "status"              : serviceContainer.status if serviceContainer else "stopped",
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
            ("Name", "Container", "Image", "Status", "IP"),
        ]
        for service in serviceData:
            tableData.append(
                (
                    serviceData[service]["name"],
                    serviceData[service]["container_name"],
                    serviceData[service]["docker_image"],
                    serviceData[service]["status"],
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
        {--c|command=? : Command to run. (Default=bash)}
        {--u|user=? : User to run command as. (Default=root)}
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