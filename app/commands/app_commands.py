from cleo import Command
from app.commands import getProject, getAppsToInvoke

class AppShell(Command):
    """
    Shell in to application container.

    app:shell
        {app? : Application to shell in to. (First available if not provided.)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        appName = self.argument("app")
        project = getProject(self)
        for app in project.getApplications():
            if app.config.getName() == appName or not appName:
                app.shell()
                break
