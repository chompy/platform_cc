import os
import json

class PlatformConfig:

    """ Global configuration. """

    def __init__(self, filename = None):
        if not filename:
            filename = os.path.join(os.path.expanduser("~"), ".platform_cc_config.json")
        self.filename = filename
        self.config = {}
        self.load()
            
    def load(self):
        """ Load config from file. """
        if os.path.exists(self.filename):
            with open(self.filename, "r") as f:
                self.config = json.load(f)

    def save(self):
        """ Save config to file. """
        with open(self.filename, "w") as f:
            json.dump(self.config, f)

    def get(self, name, default = None):
        """ Retrieve configuration value. """
        return self.config.get(name, default)

    def set(self, name, value):
        """ Set configuration value. """
        self.config[name] = value