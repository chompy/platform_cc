#!/usr/bin/env python
# -*- coding: utf-8 -*-

from cleo import Application
from app.commands.project_commands import ProjectStart, ProjectStop, ProjectBuild, ProjectDeploy
from app.commands.var_commands import VarSet, VarGet, VarDelete, VarList

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

if __name__ == '__main__':
    cleoApp.run()

"""
# display command list
if args.action == "list":
    print_help()
    sys.exit()

# get platform project
project = PlatformProject()

# get apps in project
apps = project.getApplications()
if not apps:
    log_stdout("No apps available.")
    sys.exit()

# perform action
appsArg = []
if args.apps:
    appsArg = args.apps.strip().lower().split(",")
if appsArg:
    for app in apps:
        if app.config.getName().lower() not in appsArg:
            apps.remove(app)
    if not apps:
        log_stdout("No apps available.")
        sys.exit()

if args.action == "start":
    for app in apps: app.start()
    project.router.start()
elif args.action == "stop":
    for app in apps: app.stop()
    project.router.stop()
elif args.action == "build":
    for app in apps: app.build()
elif args.action == "deploy":
    for app in apps: app.deploy()

# vars
elif args.action == "var:list":
    print project.vars.all()
elif args.action[:3] == "var":
    print varArgs
    if args.action == "var:set":
        project.vars.set(
            varArgs.key,
            varArgs.value
        )
    elif args.action == "var:get":
        print project.vars.get(varArgs.key)
    elif args.action == "var:delete":
        project.vars.delete(varArgs.key)
"""