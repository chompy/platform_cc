import os
import sys
from app.platform_project import PlatformProject
from app.platform_utils import log_stdout

# get args
args = sys.argv[1:]

print "Platform.CC (v0.01) -- By Contextual Code"

def print_help():
    print ""
    print "Usage:"
    print "  command [options] [arguments]"
    print ""
    print "Options:"
    print "  -a, --apps              Target just the specified apps (comma delimited)."
    print ""
    print "Available commands:"
    print " start                    Start application(s)."
    print " stop                     Stop application(s)."
    print " build                    Build running application(s)."
    print " var"
    print "  var:set                 Set a project variable."
    print "  var:get                 Get a project variable."
    print "  var:list                List all project variables."

# no action, display help and exit
if len(args) == 0:
    print_help()
    sys.exit()

# get action
action = args[0].strip().lower()

# 'help' action
if action in ["help", "-h", "--help"]:
    print_help()
    sys.exit()

# get platform project
project = PlatformProject()

# get apps in project
apps = project.getApplications()
if not apps:
    log_stdout("No apps available.")
    sys.exit()

# parse remaining arguments and options
OPTION_KEYS = {
    "apps" : ["-a", "--apps=", "--apps "]
}
options = {}
additionalArguments = []
if len(args) > 1:
    for arg in args[1:]:
        for key in OPTION_KEYS:
            for argOption in OPTION_KEYS[key]:
                if arg[:len(argOption)] == argOption:
                    options[key] = arg[len(argOption):]
        if arg[0] != "-":
            additionalArguments.append(arg)

# perform action
appsArg = []
if options.get("apps", []):
    appsArg = options["apps"].strip().lower().split(",")
if appsArg:
    for app in apps:
        if app.config.getName().lower() not in appsArg:
            apps.remove(app)
    if not apps:
        log_stdout("No apps available.")
        sys.exit()

if action == "start":
    for app in apps: app.start()
    project.router.start()
elif action == "stop":
    for app in apps: app.stop()
    project.router.stop()
elif action == "build":
    for app in apps: app.build()
elif action == "var:set":
    if len(additionalArguments) < 2:
        log_stdout("You must provide a key and value pair to set a variable.")
        sys.exit()
    project.vars.set(
        additionalArguments[0],
        additionalArguments[1]
    )
elif action == "var:get":
    if len(additionalArguments) < 1:
        log_stdout("You must provide a key to get a variable.")
        sys.exit()
    print project.vars.get(
        additionalArguments[0]
    )
elif action == "var:list":
    print project.vars.all()