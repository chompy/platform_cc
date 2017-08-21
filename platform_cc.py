import os
import sys
import argparse
from app.platform_app import PlatformApp

# get args
parser = argparse.ArgumentParser()
parser.add_argument(
    "action",
    type=str,
    help="Action to perform. (start|stop)"
)
args = parser.parse_args()

# get platform app
platform = PlatformApp()

# perform action
action = args.action.strip().lower()
if action == "start":
    platform.start()
elif action == "stop":
    platform.stop()
elif action == "build":
    platform.build()