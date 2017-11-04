from cleo import Command
from app.commands import getProject

class AppShell(Command):
    """
    Shell in to application container.

    app:shell
        {app? : Application to shell in to. (First available if not provided.)}
        {--c|command=? : Command to run.}
        {--u|user=? : User to shell as. (Default=web)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        appName = self.argument("app")
        command = self.option("command")
        user = self.option("user")
        project = getProject(self)
        for app in project.getApplications():
            if app.config.getName() == appName or not appName:
                app.shell(
                    command if command else "bash",
                    user if user else "web"
                )
                break
