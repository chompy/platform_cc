import collections
from urlparse import urlparse
from config.platform_router_config import PlatformRouterConfig
from platform_docker import PlatformDocker
from app.platform_utils import log_stdout

class PlatformRouter:
    """ Provide router to route request to specific app. """

    ROUTE_DOMAIN_REPLACE = "{default}"

    def __init__(self, projectHash, projectPath, projectVars, projectApps):
        self.config = PlatformRouterConfig(
            projectHash,
            projectPath
        )
        self.projectVars = projectVars
        self.projectApps = projectApps
        self.docker = PlatformDocker(self.config)

    def generateNginxConfig(self):
        """ Generate nginx config file for application. """

        serverList = collections.OrderedDict()
        routes = self.config.getRoutes()
        for routeSyntax in routes:

            parseRouteSyntax = urlparse(routeSyntax)
            isHttps = parseRouteSyntax.scheme == "https"
            serverKey = "%s:%s" % (
                "https" if isHttps else "http",
                parseRouteSyntax.hostname
            )
            if not serverKey in serverList:
                serverList[serverKey] = {
                    "scheme" :              "https" if isHttps else "http",
                    "hostname" :            parseRouteSyntax.hostname,
                    "paths" :               {}
                }

            if parseRouteSyntax.path not in serverList[serverKey]["paths"]:
                serverList[serverKey]["paths"][parseRouteSyntax.path] = {
                    "type" :                routes[routeSyntax].get("type", "upstream"),
                    "upstream" :            routes[routeSyntax].get("upstream", "",),
                    "to" :                  routes[routeSyntax].get("to", "")
                }
        
        nginxConf = ""
        for serverName in serverList:
            nginxConf += "server {\n"
            nginxConf += "\tserver_name %s;\n" % (
                serverList[serverName]["hostname"].replace(
                    self.ROUTE_DOMAIN_REPLACE,
                    "_"
                )
            )
            # TODO HTTPS
            nginxConf += "\tlisten 80;\n"

            paths = serverList[serverName]["paths"]
            for path in paths:
                nginxConf += "\tlocation %s {\n" % (
                    path
                )
                if paths[path]["type"] == "upstream":
                    upstream = paths[path]["upstream"].split(":")[0]
                    for app in self.projectApps:
                        if app.config.getName() == upstream:
                            nginxConf += "\t\tproxy_pass http://%s;\n" % (
                                app.web.docker.containerId
                            )
                            break
                nginxConf += "\t}\n"
            nginxConf += "}\n"

        routerProvisionConfig = self.docker.getProvisioner().config
        baseNginxConfig = routerProvisionConfig.get("router_conf", "")
        return baseNginxConfig.replace(
            "{{ROUTER_SERVERS}}",
            nginxConf
        )

    def start(self):
        """ Start router. """
        log_stdout("Starting router.")
        self.docker.start()
        self.docker.getProvisioner().copyStringToFile(
            self.generateNginxConfig(),
            "/etc/nginx/nginx.conf"
        )
        self.docker.getContainer().restart()

    def stop(self):
        """ Stop router. """
        log_stdout("Stopping router.")
        self.docker.stop()        