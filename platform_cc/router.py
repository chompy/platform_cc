from container import Container

class PlatformRouter(Container):

    """
    Main router for accessing all projects via the web.
    """

    def __init__(self, dockerClient = None):
        Container.__init__(self, {}, "router", dockerClient)

    def getDockerImage(self):
        return "nginx:1.13"

    def getContainerName(self):
        return "%s%s" % (
            self.CONTAINER_NAME_PREFIX,
            self.name
        )
    
    def getContainerPorts(self):
        return {
            "80/tcp"        : "80/tcp",
            "443/tcp"       : "443/tcp"
        }

    def getNetworkName(self):
        return "bridge"

    def getVolume(self, name = ""):
        return None