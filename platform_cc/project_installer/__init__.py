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

def projectInstall(project, config = {}, startFrom = 0):
    """
    Run tasks defined in config on given project.

    :param tasks: Platform.CC project
    :param config: Install tasks config
    :param startFrom: Specify task/step to start from
    """
    logger = logging.getLogger(__name__)
    config = dict(config)
    # set vars
    for key, value in config.get("vars", {}).items():
        project.variables.set(key, value)

    # install ssh keys
    sshKey = config.get("ssh_key", "")
    sshKnownHosts = config.get("ssh_known_hosts", "")
    if sshKey:
        logger.info("Install SSH key.")
        project.config.set("ssh_key", sshKey)
    if sshKnownHosts:
        logger.info("Install SSH known hosts.")
        project.config.set("ssh_known_hosts", sshKnownHosts)

    # start project
    project.start()

    # start from should be >= 1
    if not startFrom: startFrom = 1
    startFrom = int(startFrom)
    if startFrom <= 0: startFrom = 1

    # tasks
    tasks = config.get("tasks", [])
    for index in range(len(tasks)):
        taskParams = tasks[index]
        logger.info("STEP %d - Execute install task '%s.'" % (index + 1, taskParams.get("type")))
        if index + 1 < startFrom:
            logger.info("Task skipped.")
            continue
        task = getTaskHandler(project, taskParams)
        if not task.checkCondition():
            logger.info("Task run condition did not pass, skipped.")
            continue
        task.run()

    # run deploy hooks
    applications = project.dockerFetch(filterType="application")
    for application in applications:
        application.deploy()
