import yaml
import yamlordereddictloader

class BasePlatformParser:
    """
    Base class for Platform.sh configuration parser.
    """

    def __init__(self):
        pass

    def _readYaml(self, path):
        """
        Read YAML configuration file.

        :param path: Path to YAML file
        :return: Dictionary containing configuration.
        """
        config = None
        with open(path, "r") as f:
            config = yaml.load(f, Loader=yamlordereddictloader.Loader)
        return config