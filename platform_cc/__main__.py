#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import absolute_import
import sys
import os
import logging
sys.path.append(os.path.dirname(__file__))
import pkg_resources
from cleo import Application
from commands.variables import VariableSet, VariableGet, VariableDelete, VariableList
from commands.services import ServiceStart, ServiceStop, ServiceRestart, ServiceList, ServiceShell
from commands.applications import ApplicationStart

# fetch version
try:
    version = pkg_resources.require("platform_cc")[0].version
except pkg_resources.DistributionNotFound:
    version = "vDEVELOPMENT"

# init logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

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

def main():
    cleoApp.run()

if __name__ == '__main__':
    main()
