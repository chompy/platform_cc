import os
import json
import time
import hashlib
import base36
import random
from variables import getVariableStorage
from variables.json import JsonVariables

class PlatformProject:
    """
    Container class for all elements of a Platform.sh
    project. (Applications, services, variables, etc).
    """

    """ Filename of project config file. """
    PROJECT_CONFIG_FILE = ".pcc_project.json"

    """ Salt used to generate project unique ids. """
    HASH_SALT = "6fabb8b0ee9&(2cae2eb26306cdc51012f180eb$NBd!a0e"

    def __init__(self, path):
        """
        Constructor.

        :param path: Path to project root
        """

        # set project path
        self.path = str(path)

        # validate project path
        if not os.path.isdir(self.path):
            raise ValueError("Invalid project path.")

        # load config (use JsonVariables class to do this
        # as it already contains the functionality)
        self.config = JsonVariables(
            self.path,
            {
                "variables_json_filename" : self.PROJECT_CONFIG_FILE
            }
        )

        # generate uid if it does not exist
        if not self.config.get("uid"):
            self.config.set("uid", self._generateUid())

        # get variable storage
        self.variables = getVariableStorage(
            self.path,
            self.config
        )

    def _generateUid(self):
        """
        Generate a unique id for the project.

        :param path: Path to project root
        :return: Unique id string
        :rtype: str
        """
        return base36.dumps(
            int(
                hashlib.sha256(
                    (
                        "%s-%s-%s-%s" % (
                            self.HASH_SALT,
                            self.path,
                            str(random.random()),
                            str(time.time())
                        )
                    ).encode("utf-8")
                ).hexdigest(),
                16
            )
        )
