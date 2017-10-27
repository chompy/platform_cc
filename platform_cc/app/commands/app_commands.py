from cleo import Command
from app.commands import getProject, getAppsToInvoke, getLogger
from app.platform_sh_sync import PlatformShSync

class AppShell(Command):
    """
    Shell in to application container.

    app:shell
        {app? : Application to shell in to. (First available if not provided.)}
        {--c|command=? : Command to run.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        appName = self.argument("app")
        command = self.option("command")
        project = getProject(self)
        for app in project.getApplications():
            if app.config.getName() == appName or not appName:
                app.shell( command if command else "bash" )
                break

class AppPlatformShImport(Command):
    """
    Import application data from Platform.sh.

    app:platform_sh_import
        {ssh : SSH access url.}
        {app? : Application to shell in to. (First available if not provided.)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        appName = self.argument("app")
        project = getProject(self)
        for app in project.getApplications():
            if app.config.getName() == appName or not appName:

                pshSync = PlatformShSync(
                    app,
                    self.argument("ssh"),
                    getLogger(self)
                )
                print pshSync.rsyncCopy("/app/var", ["png", "gif", "jpg", "jpeg"])

                break