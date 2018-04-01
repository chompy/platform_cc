import os
from cleo import Command
from commands import getProject, outputJson, outputTable

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
