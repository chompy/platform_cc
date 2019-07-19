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

import os
import io
import tarfile
from .base import BaseTaskHandler
from ...commands.mysql import getMysqlService
from ...exception.state_error import StateError

MYSQL_IMPORT_VALID_SERVICE_TYPES = ["mysql", "mariadb"]

class MysqlImportTaskHandler(BaseTaskHandler):

    """
    Task handler for importing MySQL dumps.
    """

    @classmethod
    def getType(cls):
        return "mysql_import"

    def run(self):
        # validate params
        self.checkParams(["from", "to"])

        # parse 'to' parameter
        toParams = self.params.get("to", "").split(":")
        if len(toParams) < 2:
            raise ValueError("'mysql_import' install task requires 'to' parameter in format <service_name>:<database_name>.")
        serviceName = toParams[0].strip()
        database = toParams[1].strip()

        # get mysql service
        service = getMysqlService(self.project, serviceName)
        if not service.isRunning():
            raise StateError(
                "Service '%s' is not running." % service.getName()
            )

        # get from container and path
        fromApp, fromPath = self.parseAppPath(self.params.get("from"))

        # download dump from container
        tarStream, _ = fromApp.getContainer().get_archive(
            fromPath
        )
        tarFileObject = io.BytesIO()
        for d in tarStream:
            tarFileObject.write(d)
        tarStream.close()

        # build command to run
        cmd = "mysql -h 127.0.0.1 -uroot --password=\"%s\"" % (
            service.getPassword()
        )
        if self.params.get("to"):
            cmd += " --database=\"%s\"" % database
        
        # upload dump
        fileBaseName = os.path.basename(fromPath)
        fileExt = os.path.splitext(fileBaseName)[1]
        service.uploadFile(tarFileObject, "/tmp/dump%s.tar" % fileExt)
        service.runCommand(
            "cd /tmp && tar -xf dump%s.tar && rm dump%s.tar" % (
                fileExt, fileExt
            )
        )
        tarFileObject.close()

        # un-gunzip file if gz is file extension
        if fileExt == ".gz":
            service.runCommand(
                "cd /tmp && gunzip -f %s" % fileBaseName
            )
            fileBaseName = fileBaseName[0:-3]
        
        # run command
        cmd = "sh -c 'cat /tmp/%s | %s && rm /tmp/%s'" % (fileBaseName, cmd, fileBaseName)
        service.runCommand(cmd)

        # delete dump
        if self.params.get("delete_dump"):
            fromApp.runCommand(
                "rm -f %s" % fromPath
            )