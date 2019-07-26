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
from cleo import Command
from platform_cc.commands import getProject, outputJson, outputTable
from platform_cc.config import PlatformConfig
from platform_cc.exception.project_init_error import ProjectInitError

GLOBAL_VAR_CONFIG_PREFIX = "global_project_var:"

class VariableSet(Command):
    """
    Set a project variable.

    variable:set
        {key : Name of variable to set.}
        {value : Value of variable.}
        {--p|path=? : Path to project root. (Default=current directory)}
        {--g|global : Set as global project variable. }
    """

    def handle(self):
        # set global
        if self.option("global"):
            config = PlatformConfig()
            config.set(
                "%s%s" % (GLOBAL_VAR_CONFIG_PREFIX, self.argument("key")),
                self.argument("value")
            )
            return
        # set project
        project = getProject(self)
        project.variables.set(
            self.argument("key"),
            self.argument("value")
        )

class VariableGet(Command):
    """
    Get a project variable.

    variable:get
        {key : Name of variable to get.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        self.line(
            project.variables.get(
                self.argument("key")
            )
        )

class VariableDelete(Command):
    """
    Delete a project variable.

    variable:delete
        {key : Name of variable to delete.}
        {--p|path=? : Path to project root. (Default=current directory)}
        {--g|global : Delete global project variable. }
    """

    def handle(self):
        # delete global var
        if self.option("global"):
            config = PlatformConfig()
            config.delete(
                "%s%s" % (GLOBAL_VAR_CONFIG_PREFIX, self.argument("key"))
            )
            return
        # delete project var
        project = getProject(self)
        project.variables.delete(
            self.argument("key")
        )

class VariableList(Command):
    """
    List all project variables.

    variable:list
        {--p|path=? : Path to project root. (Default=current directory)}
        {--j|json : Output as JSON.}
        {--g|global : Only display global project variable. }
    """

    def handle(self):
        allVars = None
        project = None
        if not self.option("global"):
            try:
                project = getProject(self)
                allVars = project.variables.all()
            except ProjectInitError:
                pass
        if not allVars:
            allVars = PlatformConfig().getGlobalProjectVars()

        # json output
        if self.option("json"):
            outputJson(self, allVars)
            return

        # terminal tables output
        tableData = [
            ("Key", "Value")
        ]
        for key in allVars:
            tableData.append(
                (key, allVars[key])
            )
        outputTable(
            self,
            "Project '%s' - Variables" % project.getUid()[0:6] if project else "Global Project Variables",
            tableData
        )
