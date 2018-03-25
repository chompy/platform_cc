import os
from cleo import Command
from project import PlatformProject

class VariableSet(Command):
    """
    Set a project variable.

    var:set
        {key : Name of variable to set.}
        {value : Value of variable.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        path = self.option("path")
        if not path: path = os.getcwd()
        project = PlatformProject(path)
        project.variables.set(
            self.argument("key"),
            self.argument("value")
        )

class VariableGet(Command):
    """
    Get a project variable.

    var:get
        {key : Name of variable to set.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        path = self.option("path")
        if not path: path = os.getcwd()
        project = PlatformProject(path)
        print(
            project.variables.get(
                self.argument("key")
            )
        )

class VariableDelete(Command):
    """
    Delete a project variable.

    var:delete
        {key : Name of variable to set.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        path = self.option("path")
        if not path: path = os.getcwd()
        project = PlatformProject(path)
        project.variables.delete(
            self.argument("key")
        )

class VariableList(Command):
    """
    List all project variables as JSON.

    var:list
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        path = self.option("path")
        if not path: path = os.getcwd()
        project = PlatformProject(path)
        print(
            project.variables.all()
        )
