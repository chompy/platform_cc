from .php import PhpApplication

""" Map application names to their application class. """
APPLICATION_MAP = {
    "php"           : PhpApplication
}

def getApplication(project, config):
    """
    Get application handler from configuration.
    
    :param project: Dictionary containing project data
    :param config: Dictionary containing application configuration
    :return: Application handler object
    :rtype: .base.BasePlatformApplication
    """

    # validate config
    if not isinstance(config, dict):
        raise ValueError("Config parameter must be a dictionary (dict) object.")
    if "type" not in config or "name" not in config:
        raise ValueError("Config parameter is missing parameters required for an application.")
    appType = config["type"].split(":")[0]
    if appType not in APPLICATION_MAP:
        raise NotImplementedError(
            "No appliocation handler available for '%s.'" % (
                config["type"]
            )
        )
    # init application
    return APPLICATION_MAP[appType](
        project,
        config
    )