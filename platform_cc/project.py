import os
import json
import time
import hashlib
import base36
import random
from variables import getVariableStorage

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

        # load config
        self.config = {}
        self._loadConfig()

        # generate uid if it does not exist
        if "uid" not in self.config:
            self.config["uid"] = self._generateUid()
            self._saveConfig()

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

    def _loadConfig(self):
        """
        Open config json file and overwrite config currently
        in memory.
        """
        # validate path
        if not os.path.isdir(self.path):
            raise ValueError("Project path '%s' is invalid." % self.path)
        # get path to config
        configPath = os.path.join(self.path, self.PROJECT_CONFIG_FILE)
        # open+parse config
        if os.path.isfile(configPath):
            with open(configPath, "r") as f:
                self.config = json.load(f)

    def _saveConfig(self):
        """
        Open config json file to store config currently in
        memory.
        """
        # validate path
        if not os.path.isdir(self.path):
            raise ValueError("Project path '%s' is invalid." % self.path)
        # get path to config
        configPath = os.path.join(self.path, self.PROJECT_CONFIG_FILE)
        # open+write to config
        with open(configPath, "w") as f:
            json.dump(
                self.config,
                f,
                sort_keys=True,
                indent=4,
                separators=(',', ': ')
            )