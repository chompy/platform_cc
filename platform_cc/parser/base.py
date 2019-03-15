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
import yaml
import yamlordereddictloader
from .yaml_archive import ArchiveTag


class BasePlatformParser:
    """
    Base class for Platform.sh configuration parser.
    """

    def __init__(self, projectPath):
        """
        Constructor.

        :param projectPath: Path to project root
        """

        yaml.SafeLoader.add_constructor(
            ArchiveTag.yaml_tag, ArchiveTag.from_yaml
        )
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

    def _readYaml(self, path):
        """
        Read single YAML configuration.

        :param path: Path string to YAML file
        :return: Dictionary containing configuration.
        """
        ArchiveTag.base_path = os.path.dirname(path)
        loadConf = {}
        with open(path, "r") as f:
            loadConf = yaml.load(f, Loader=yamlordereddictloader.SafeLoader)
        return loadConf

    def _readYamls(self, paths):
        """
        Read YAML configuration files and merge them all together
        in to a single configuration dict.

        :param paths: Path or list of paths to YAML file
        :return: Dictionary containing configuration.
        """
        config = {}
        if type(paths) is not list:
            paths = [paths]
        for path in paths:
            self._mergeDict(
                self._readYaml(path),
                config
            )
        return config
