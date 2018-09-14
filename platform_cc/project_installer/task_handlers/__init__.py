from .asset_s3 import AssetS3TaskHandler

TASK_HANDLER_MAP = {
    "asset_s3"      : AssetS3TaskHandler
}

def getTaskHandler(project, params = {}):
    """
    Get task handler by its type name.

    :param project: Project
    :param params: Dict containing task parameters:
    :rtype: .base.BaseTaskHandler
    """
    type = dict(params).get("type")
    if not type:
        raise ValueError("Task handler parameters much contain a 'type.'")
    taskHandler = TASK_HANDLER_MAP.get(type)
    if not taskHandler:
        raise NotImplementedError("Task handler '%s' has not been implemented." % type)
    return taskHandler(project, params)