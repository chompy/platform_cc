import os
import yaml
import logging
from .task_handlers import getTaskHandler

def loadInstallFile(path):
    logger = logging.getLogger(__name__)
    logger.info("Load project installer YAML file from '%s.'" % path)
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
    logger = logging.getLogger(__name__)
    config = dict(config)
    # set vars
    for key, value in config.get("vars", {}).items():
        project.variables.set(key, value)

    # start project
    project.start()

    # tasks
    for taskParams in config.get("tasks", []):
        logger.info("Execute install task '%s.'" % taskParams.get("type"))
        task = getTaskHandler(project, taskParams)
        if not task.checkCondition():
            logger.info("Task run condition did not pass, skipped.")
            continue
        task.run()

    # run deploy hooks
    applications = project.dockerFetch(filterType="application")
    for application in applications:
        application.deploy()
