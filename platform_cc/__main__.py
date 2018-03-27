#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import absolute_import
import sys
import os
sys.path.append(os.path.dirname(__file__))
import pkg_resources
from cleo import Application
from commands.variables import VariableSet, VariableGet, VariableDelete, VariableList
from commands.services import ServiceStart, ServiceStop, ServiceRestart, ServiceList, ServiceShell

try:
    version = pkg_resources.require("platform_cc")[0].version
except pkg_resources.DistributionNotFound:
    version = "vDEVELOPMENT"

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

def main():
    cleoApp.run()

if __name__ == '__main__':
    main()
