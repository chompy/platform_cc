import os
from cleo import Command
from platform_cc.commands import getProject, outputJson, outputTable

class ApplicationStart(Command):
    """
    Start one or more applications.

    application:start
        {name* : Name(s) of application.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        for name in self.argument("name"):
            application = project.getApplication(name)
            application.start()
            self.line(application.getContainerName())

class ApplicationStop(Command):
    """
    Stop one or more applications.

    application:stop
        {name* : Name(s) of application.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        for name in self.argument("name"):
            application = project.getApplication(name)
            application.stop()
            self.line(application.getContainerName())

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
                "docker_image"        : application.getDockerImage(),
                "status"              : appContainer.status if appContainer else "stopped",
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
            ("Name", "Container", "Image", "Status", "IP"),
        ]
        for application in appData:
            tableData.append(
                (
                    appData[application]["name"],
                    appData[application]["container_name"],
                    appData[application]["docker_image"],
                    appData[application]["status"],
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
        {--c|command=? : Command to run. (Default=bash)}
        {--u|user=? : User to run command as. (Default=web)}
    """

    def handle(self):
        project = getProject(self)
        name = self.option("name")
        if not name:
            name = project.getApplicationsParser().getApplicationNames()[0]
        application = project.getApplication(name)
        command = self.option("command")
        if not command: command = "bash"
        user = self.option("user")
        if not user: user = "web"
        application.shell(
            command,
            user
        )