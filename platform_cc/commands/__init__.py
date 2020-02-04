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
import sys
import os
import logging.config
import json
from cleo import Application
from .variables import VariableSet, VariableGet, VariableDelete, VariableList
from .services import ServiceRestart, ServiceList, ServiceShell, ServicePull
from .applications import ApplicationRestart, ApplicationList, ApplicationShell, ApplicationBuild, ApplicationDeployHook, ApplicationPull, ApplicationNginxConfig
from .router import RouterStart, RouterStop, RouterRestart, RouterAdd, RouterRemove, RouterNginx
from .project import ProjectStart, ProjectStop, ProjectRestart, ProjectRoutes, ProjectOptionSet, ProjectOptionList, ProjectPurge, ProjectInstall, ProjectPull
from .mysql import MysqlSql, MysqlDump
from .all import AllStop, AllPurge, AllList
from .platformsh import PlatformShLogin, PlatformShLogout, PlatformShClone, PlatformShSetSsh, PlatformShSync
from ..core.version import PCC_VERSION
from ..core import LOGGING_CONFIG_JSON

loggingConfig = {}
with open(LOGGING_CONFIG_JSON, "rt") as f:
    loggingConfig = json.load(f)
logging.config.dictConfig(loggingConfig)

# init cleo
cleoApp = Application(
    "Platform.CC -- By Contextual Code",
    PCC_VERSION
)
cleoApp.add(VariableSet())
cleoApp.add(VariableGet())
cleoApp.add(VariableDelete())
cleoApp.add(VariableList())
cleoApp.add(ServiceRestart())
cleoApp.add(ServiceList())
cleoApp.add(ServiceShell())
cleoApp.add(ServicePull())
cleoApp.add(ApplicationRestart())
cleoApp.add(ApplicationList())
cleoApp.add(ApplicationShell())
cleoApp.add(ApplicationBuild())
cleoApp.add(ApplicationDeployHook())
cleoApp.add(ApplicationPull())
cleoApp.add(ApplicationNginxConfig())
cleoApp.add(RouterStart())
cleoApp.add(RouterStop())
cleoApp.add(RouterRestart())
cleoApp.add(RouterAdd())
cleoApp.add(RouterRemove())
cleoApp.add(RouterNginx())
cleoApp.add(ProjectStart())
cleoApp.add(ProjectStop())
cleoApp.add(ProjectRestart())
cleoApp.add(ProjectRoutes())
cleoApp.add(ProjectOptionSet())
cleoApp.add(ProjectOptionList())
cleoApp.add(ProjectPurge())
cleoApp.add(ProjectInstall())
cleoApp.add(ProjectPull())
cleoApp.add(MysqlSql())
cleoApp.add(MysqlDump())
cleoApp.add(AllStop())
cleoApp.add(AllPurge())
cleoApp.add(AllList())
cleoApp.add(PlatformShLogin())
cleoApp.add(PlatformShLogout())
cleoApp.add(PlatformShClone())
cleoApp.add(PlatformShSetSsh())
cleoApp.add(PlatformShSync())

def main():
    cleoApp.run()
