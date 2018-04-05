#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import absolute_import
import sys
import os
import logging.config
import json
sys.path.append(os.path.dirname(__file__))
import pkg_resources
from cleo import Application
from commands.variables import VariableSet, VariableGet, VariableDelete, VariableList
from commands.services import ServiceStart, ServiceStop, ServiceRestart, ServiceList, ServiceShell
from commands.applications import ApplicationStart, ApplicationStop, ApplicationList, ApplicationShell
from commands.router import RouterStart, RouterStop, RouterAdd, RouterRemove
from commands.project import ProjectStart, ProjectStop
from commands.mysql import MysqlSql

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
cleoApp.add(ApplicationList())
cleoApp.add(ApplicationShell())
cleoApp.add(RouterStart())
cleoApp.add(RouterStop())
cleoApp.add(RouterAdd())
cleoApp.add(RouterRemove())
cleoApp.add(ProjectStart())
cleoApp.add(ProjectStop())
cleoApp.add(MysqlSql())

def main():
    cleoApp.run()

if __name__ == '__main__':
    main()
