"""
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
"""

from .asset_s3 import AssetS3TaskHandler
from .mysql_import import MysqlImportTaskHandler
from .shell import ShellTaskHandler
from .rsync import RsyncTaskHandler

TASK_HANDLER_MAP = {
    "asset_s3"      : AssetS3TaskHandler,
    "mysql_import"  : MysqlImportTaskHandler,
    "shell"         : ShellTaskHandler,
    "rsync"         : RsyncTaskHandler
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