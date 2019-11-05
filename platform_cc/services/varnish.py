from .base import BasePlatformService
import io

class VarnishService(BasePlatformService):
    """
    Handler for Varnish service.
    Just a placeholder/dummy service for now.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "varnish:5.2":            "plopix/docker-varnish5",
        "varnish:6.0":            "plopix/docker-varnish6"
    }

    def getStartGroup(self):
        return self.START_POST_APP_A

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerVolumes(self):
        return {}

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"][self.getName()] = {
            "host":           self.getContainerName(),
            "ip":             data.get("ip", ""),
            "scheme":         "http",
            "port":           80
        }
        data["platform_relationships"]["varnish"] = (
            data["platform_relationships"][self.getName()]
        )
        return data

    def generateVcl(self, hostmap = {}):
        vclStr = "vcl 4.1;\n"
        vclStr += "import std;\n"
        vclStr += "import directors;\n"
        for name, value in self.config.get("_relationships", {}).items():
            value = value.strip().split(":")
            if value[1] != "http": continue
            hostname = "pcc_%s_%s" % (self.project.get("short_uid"), value[0])
            if name in hostmap:
                hostname = hostmap[name]
            vclStr += "backend %s-server {\n    .host = \"%s\";\n    .port = \"%s\";\n}\n" % (
                name, hostname, 80
            )
        vclStr += "sub vcl_init {\n"
        for name, value in self.config.get("_relationships", {}).items():
            vclStr += "new %s = directors.round_robin();\n%s.add_backend(%s-server);\n" % (
                name,
                name,
                name
            )
        vclStr += "}\n"
        vclStr += self.config.get("vcl", "")
        return vclStr

    def start(self):
        BasePlatformService.start(self)
        # add vcl
        fObj = io.BytesIO(self.generateVcl().encode())
        self.uploadFile(
            fObj,
            "/etc/varnish/default.vcl"
        )
        # restart for config to take effect
        self.getContainer().restart()