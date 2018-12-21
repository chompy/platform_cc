import os
import io
import tarfile
from .base import BaseTaskHandler
from ...exception.state_error import StateError

class ShellTaskHandler(BaseTaskHandler):

    """
    Task handler for running shell commands.
    """

    @classmethod
    def getType(cls):
        return "shell"

    def run(self):
        # validate params
        self.checkParams(["to", "command"])
        # get 'to' application
        toApp = self.project.getApplication(self.params.get("to"))
        # get user to run as
        user = self.params.get("user", "web")
        # run command
        toApp.runCommand(self.params.get("command"), user)
