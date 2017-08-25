from config.platform_router_config import PlatformRouterConfig
from platform_docker import PlatformDocker
from app.platform_utils import log_stdout

class PlatformRouter:
    """ Provide router to route request to specific app. """

    def __init__(self, projectHash, projectPath):
        self.config = PlatformRouterConfig(
            projectHash,
            projectPath
        )
        self.docker = PlatformDocker(self.config)

    def generateNginxConfig(self):
        """ Generate nginx config file for application. """

        routes = self.config.getRoutes()
        for routeSyntax in routes:

            isHttps = True if routeSyntax[:5].lower() == "https" else False
            routeType = routes[routeSyntax].get("type", "upstream")
            upstream = routes[routeSyntax].get("upstream", "")
            

            print routeSyntax

        return "server { }"


    def start(self):
        """ Start router. """
        log_stdout("Starting router.")
        print self.generateNginxConfig()
        self.docker.start()
        #self.docker.getProvisioner().copyStringToFile(
        #    self.generateNginxConfig(),
        #    "/etc/nginx/nginx.conf"
        #)
        #self.docker.getContainer().restart()

    def stop(self):
        """ Stop router. """
        log_stdout("Stopping router.")
        self.docker.stop()        