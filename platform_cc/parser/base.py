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

    def _mergeDict(self, source, dest):
        """
        Utility method, merge two dicts recursively.
        See... https://stackoverflow.com/a/20666342

        :param source: source dict
        :param dest: destination dict
        """
        for key, value in source.items():
            if isinstance(value, dict):
                node = dest.setdefault(key, {})
                self._mergeDict(value, node)
            else:
                dest[key] = value
        return dest

    def _readYaml(self, paths):
        """
        Read YAML configuration file.

        :param paths: Path or list of paths to YAML file
        :return: Dictionary containing configuration.
        """
        config = {}
        if type(paths) is not list:
            paths = [paths]
        for path in paths:
            with open(path, "r") as f:
                self._mergeDict(
                    yaml.load(f, Loader=yamlordereddictloader.Loader),
                    config
                )
        return config