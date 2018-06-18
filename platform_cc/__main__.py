#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import absolute_import
import sys
import os
import logging.config
import json
import pkg_resources
from cleo import Application
from platform_cc.commands.variables import VariableSet, VariableGet, VariableDelete, VariableList
from platform_cc.commands.services import ServiceStart, ServiceStop, ServiceRestart, ServiceList, ServiceShell
from platform_cc.commands.applications import ApplicationStart, ApplicationStop, ApplicationRestart, ApplicationList, ApplicationShell, ApplicationBuild, ApplicationDeployHook, ApplicationPull
from platform_cc.commands.router import RouterStart, RouterStop, RouterRestart, RouterAdd, RouterRemove
from platform_cc.commands.project import ProjectStart, ProjectStop, ProjectRestart, ProjectRoutes, ProjectOptionSet, ProjectOptionList, ProjectPurge
from platform_cc.commands.mysql import MysqlSql

# fetch version
try:
    version = pkg_resources.require("platform_cc")[0].version
except pkg_resources.DistributionNotFound:
    version = "vDEVELOPMENT"

# init logging
LOGGING_CONFIG_JSON = os.path.join(
    os.path.dirname(__file__),
    "logging.json"
)
loggingConfig = {}
with open(LOGGING_CONFIG_JSON, "rt") as f:
    loggingConfig = json.load(f)
logging.config.dictConfig(loggingConfig)

# init cleo
cleoApp = Application(
    "Platform.CC -- By Contextual Code",
    version
)
cleoApp.add(VariableSet())
cleoApp.add(VariableGet())
cleoApp.add(VariableDelete())
cleoApp.add(VariableList())
cleoApp.add(ServiceStart())
cleoApp.add(ServiceStop())
cleoApp.add(ServiceRestart())
cleoApp.add(ServiceList())
cleoApp.add(ServiceShell())
cleoApp.add(ApplicationStart())
cleoApp.add(ApplicationStop())
cleoApp.add(ApplicationRestart())
cleoApp.add(ApplicationList())
cleoApp.add(ApplicationShell())
cleoApp.add(ApplicationBuild())
cleoApp.add(ApplicationDeployHook())
cleoApp.add(ApplicationPull())
cleoApp.add(RouterStart())
cleoApp.add(RouterStop())
cleoApp.add(RouterRestart())
cleoApp.add(RouterAdd())
cleoApp.add(RouterRemove())
cleoApp.add(ProjectStart())
cleoApp.add(ProjectStop())
cleoApp.add(ProjectRestart())
cleoApp.add(ProjectRoutes())
cleoApp.add(ProjectOptionSet())
cleoApp.add(ProjectOptionList())
cleoApp.add(ProjectPurge())
cleoApp.add(MysqlSql())

def main():
    cleoApp.run()

if __name__ == '__main__':
    main()
