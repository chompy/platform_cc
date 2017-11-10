#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import absolute_import
import sys
import os
sys.path.append(os.path.dirname(__file__))
import pkg_resources
from cleo import Application
from app.commands.project_commands import ProjectStart, ProjectStop, ProjectProvision, ProjectDeploy, ProjectInfo, ProjectPurge
from app.commands.router_commands import RouterStart, RouterStop, RouterAdd, RouterRemove
from app.commands.var_commands import VarSet, VarGet, VarDelete, VarList
from app.commands.app_commands import AppShell
from app.commands.mysql_commands import MysqlSql, MysqlDump

try:
    version = pkg_resources.require("platform_cc")[0].version
except pkg_resources.DistributionNotFound:
    version = "vDEVELOPMENT"

cleoApp = Application(
    "Platform.CC -- By Contextual Code",
    version
)
cleoApp.add(ProjectStart())
cleoApp.add(ProjectStop())
cleoApp.add(ProjectProvision())
cleoApp.add(ProjectDeploy())
cleoApp.add(ProjectInfo())
cleoApp.add(ProjectPurge())
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
