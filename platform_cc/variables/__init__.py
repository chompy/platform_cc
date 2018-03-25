from .json import JsonVariables

def getVariableStorage(projectPath, projectConfig = {}):
    """
    Init a variable storage class from project configuration.
    """

    # default is json
    variableStorageName = "json"

    # retrieve storage method from project config
    if "variable_storage" in projectConfig:
        variableStorageName = str(projectConfig["variable_storage"])

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
    