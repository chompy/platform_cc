from .json import JsonVariables

def getVariableStorage(projectPath, projectConfig = None):
    """
    Init a variable storage class from project configuration.
    """

    # default is json
    variableStorageName = "json"

    # retrieve storage method from project config
    variableStorageName = str(
        projectConfig.get("variable_storage", variableStorageName)
    )
    
    # json
    if variableStorageName == "json":
        return JsonVariables(
            projectPath,
            projectConfig
        )

    # not found
    raise NotImplementedError(
        "Variable storage handler '%s' has not been implemented." % (
            variableStorageName
        )
    )
    