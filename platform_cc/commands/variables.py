from cleo import Command
from commands import getProject

class VariableSet(Command):
    """
    Set a project variable.

    variable:set
        {key : Name of variable to set.}
        {value : Value of variable.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        project.variables.set(
            self.argument("key"),
            self.argument("value")
        )

class VariableGet(Command):
    """
    Get a project variable.

    variable:get
        {key : Name of variable to set.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        print(
            project.variables.get(
                self.argument("key")
            )
        )

class VariableDelete(Command):
    """
    Delete a project variable.

    variable:delete
        {key : Name of variable to set.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        project.variables.delete(
            self.argument("key")
        )

class VariableList(Command):
    """
    List all project variables as JSON.

    variable:list
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        print(
            project.variables.all()
        )
