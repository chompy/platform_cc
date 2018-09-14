import os
import yaml
from .task_handlers import getTaskHandler

def loadInstallFile(path):
    if not os.path.exists(path):
        raise ValueError("Install YAML file not found.")
    with open(path, "r") as f:
        conf = yaml.load(f)
    return conf

def projectInstall(project, config = {}):
    """
    Run tasks defined in config on given project.

    :param tasks: Platform.CC project
    :param config: Install tasks config
    """
    config = dict(config)
    # set vars
    for key, value in config.get("vars", {}).items():
        project.variables.set(key, value)

    # start project
    project.start()

    # tasks
    for taskParams in config.get("tasks", []):
        getTaskHandler(project, taskParams).run()

    # stop project
    project.stop()