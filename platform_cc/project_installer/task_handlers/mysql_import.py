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
        self.checkParams(["from"])

        # get mysql service
        service = getMysqlService(self.project, self.params.get("service"))
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
            cmd += " --database=\"%s\"" % self.params.get("to")
        
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
                "cd /tmp && gunzip -f dump%s" % fileExt
            )
        
        # run command
        cmd = "sh -c 'cat /tmp/dump | %s && rm /tmp/dump'" % cmd
        service.runCommand(cmd)

        # delete dump
        if self.params.get("delete_dump"):
            os.remove(self.params.get("from"))