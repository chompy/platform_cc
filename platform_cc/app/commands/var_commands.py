from __future__ import absolute_import
from __future__ import print_function
import json
from cleo import Command
from . import getProject
from app.platform_project import PlatformProject

class VarSet(Command):
    """
    Set a project variable.

    var:set
        {key : Name of variable to set.}
        {value : Value of variable.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        getProject(self).vars.set(
            self.argument("key"),
            self.argument("value")
        )


class VarGet(Command):
    """
    Retrieve a project variable.

    var:get
        {key : Name of variable to retrieve.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        print(
            getProject(self).vars.get(
                self.argument("key")
            )
        )

class VarDelete(Command):
    """
    Delete a project variable.

    var:delete
        {key : Name of variable to delete.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        getProject(self).vars.delete(
            self.argument("key"),
        )

class VarList(Command):
    """
    List all project variables in JSON format.

    var:list
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        print(json.dumps(getProject(self).vars.all()))