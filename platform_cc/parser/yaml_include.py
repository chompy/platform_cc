import os
import yaml
import json


class IncludeTag(yaml.YAMLObject):
    """
    YAML custom tag that packs up the contents of a file.
    """

    yaml_tag = u'!include'
    base_path = "/"

    def __init__(self, value):
        dataType = ""
        path = ""
        if type(value) is str:
            dataType = "string"
            path = value
        elif type(value) is list:
            data = {}
            for _val in value:
                data[_val[0].value] = _val[1].value
            dataType = data.get("type", "")
            path = data.get("path", "")
        
        if dataType != "string":
            raise NotImplementedError("Not able to handle YAML !include tag with given value")
        self.path = os.path.join(IncludeTag.base_path, path)

    def build(self):
        content = ""
        with open(self.path, "r") as f:
            content = str(f.read())
        return content

    @classmethod
    def from_yaml(cls, loader, node):
        return IncludeTag(node.value).build()
