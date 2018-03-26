import yaml
import yamlordereddictloader

class BasePlatformParser:
    """
    Base class for Platform.sh configuration parser.
    """

    def __init__(self, projectPath):
        """
        Constructor.

        :param projectPath: Path to project root
        """
        self.projectPath = str(projectPath)

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