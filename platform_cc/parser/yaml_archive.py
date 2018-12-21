import os
import yaml
import json

class ArchiveTag(yaml.YAMLObject):
    
    """
    YAML custom tag that packs up the contents of a directory
    in to a dictionary.
    """

    yaml_tag = u'!archive'
    base_path = "/"

    def __init__(self, path):
        self.path = os.path.join(ArchiveTag.base_path, path)

    def build(self):
        res = {}
        for dirName, _, fileList in os.walk(self.path):
            for filename in fileList:
                fullPath = os.path.join(self.path, dirName, filename)
                if not os.path.isfile(fullPath): continue
                with open(fullPath, "r") as f:
                    res[os.path.join(dirName, filename)] = str(f.read())
        return res
        
    @classmethod
    def from_yaml(cls, loader, node):
        return ArchiveTag(node.value).build()

