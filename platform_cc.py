#!/usr/bin/env python
# -*- coding: utf-8 -*-

from cleo import Application
from app.commands.project_commands import ProjectStart, ProjectStop, ProjectBuild, ProjectDeploy
from app.commands.var_commands import VarSet, VarGet, VarDelete, VarList
from app.commands.app_commands import AppShell

cleoApp = Application(
    "Platform.CC -- By Contextual Code",
    "0.01"
)
cleoApp.add(ProjectStart())
cleoApp.add(ProjectStop())
cleoApp.add(ProjectBuild())
cleoApp.add(ProjectDeploy())
cleoApp.add(VarSet())
cleoApp.add(VarGet())
cleoApp.add(VarDelete())
cleoApp.add(VarList())
cleoApp.add(AppShell())

if __name__ == '__main__':
    cleoApp.run()
