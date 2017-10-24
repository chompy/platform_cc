#!/usr/bin/env python
# -*- coding: utf-8 -*-

from cleo import Application
from app.commands.project_commands import ProjectStart, ProjectStop, ProjectBuild, ProjectDeploy, ProjectInfo
from app.commands.router_commands import RouterStart, RouterStop, RouterAdd, RouterRemove
from app.commands.var_commands import VarSet, VarGet, VarDelete, VarList
from app.commands.app_commands import AppShell
from app.commands.mysql_commands import MysqlSql, MysqlDump

cleoApp = Application(
    "Platform.CC -- By Contextual Code",
    "0.01"
)
cleoApp.add(ProjectStart())
cleoApp.add(ProjectStop())
cleoApp.add(ProjectBuild())
cleoApp.add(ProjectDeploy())
cleoApp.add(ProjectInfo())
cleoApp.add(RouterStart())
cleoApp.add(RouterStop())
cleoApp.add(RouterAdd())
cleoApp.add(RouterRemove())
cleoApp.add(VarSet())
cleoApp.add(VarGet())
cleoApp.add(VarDelete())
cleoApp.add(VarList())
cleoApp.add(AppShell())
cleoApp.add(MysqlSql())
cleoApp.add(MysqlDump())

def main():
    cleoApp.run()

if __name__ == '__main__':
    main()
