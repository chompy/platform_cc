import os
import io
import tarfile
from .base import BaseTaskHandler
from ...exception.state_error import StateError

class RsyncTaskHandler(BaseTaskHandler):

    """
    Task handler for rsyncing files.
    """

    @classmethod
    def getType(cls):
        return "rsync"

    def run(self):
        # validate params
        self.checkParams(["from", "to"])

        # parse to path
        app, appPath = self.parseAppPath(self.params.get("to"))

        # get from container and path
        fromPath = self.params.get("from")

        # add host to known_hosts
        hostSplit = fromPath.split(":")[0].split("@")
        if len(hostSplit) > 1:
            app.runCommand(
                "ssh-keyscan %s >> ~/.ssh/known_hosts" % hostSplit[1],
                user="web"
            )

        # build command
        cmd = "rsync -a"

        # private key
        privateKey = self.params.get("private_key", "")
        if privateKey:
            app.runCommand(
                "chmod 0600 %s" % privateKey
            )
            cmd += " -e \"ssh -i %s\"" % privateKey

        # includes
        includes = self.params.get("includes", [])
        if includes:
            for include in includes:
                cmd += " --include=\"%s\"" % include
        
        # excludes
        excludes = self.params.get("excludes", [])
        if excludes:
            for excude in excludes:
                cmd += " --exclude=\"%s\"" % excude

        # add paths
        cmd += " %s %s" % (
            fromPath, appPath
        )

        # run command
        app.runCommand(cmd, user="web")