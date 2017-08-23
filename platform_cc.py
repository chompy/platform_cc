import os
import sys
import argparse
from app.platform_project import PlatformProject
from app.platform_utils import log_stdout

# get args
parser = argparse.ArgumentParser()
parser.add_argument(
    "action",
    type=str,
    help="Action to perform. (start|stop)"
)
parser.add_argument(
    "--apps",
    type=str,
    help="Applications to invoke. (Comma delimited.)"
)
args = parser.parse_args()

# get platform project
project = PlatformProject()
apps = project.getApplications()
if not apps:
    log_stdout("No apps available.")
    sys.exit()

# perform action
action = args.action.strip().lower()
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

if action == "start":
    for app in apps: app.start()
elif action == "stop":
    for app in apps: app.stop()
elif action == "build":
    for app in apps: app.build()