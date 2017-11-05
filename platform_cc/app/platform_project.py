import os
import time
import hashlib
import base36
import yaml
import yamlordereddictloader
import collections
import docker
import time
from urlparse import urlparse
from Crypto.PublicKey import RSA
from terminaltables import AsciiTable
from config.platform_service_config import PlatformServiceConfig
from platform_service import PlatformService
from platform_app import PlatformApp
from config.platform_app_config import PlatformAppConfig
from platform_vars import PlatformVars

class ProjectNotFoundException(Exception):
    pass

class PlatformProject:

    """ Base class for project. """

    HASH_SECRET = "4bcc181ab1f9fcc64a8c935686b55ca794e76d63"
    PLATFORM_ROUTES_PATH = ".platform/routes.yaml"
    ROUTE_DOMAIN_REPLACE = "{default}"
    ROUTE_INTERNAL_DOMAIN_DOT_REPLACE = "-"

    def __init__(self, projectPath = "", logger = None):
        self.projectPath = projectPath
        self.logger = logger
        if not os.path.isdir(os.path.realpath(projectPath)):
            raise ProjectNotFoundException("Could not find project at '%s.'" % os.path.realpath(projectPath))
        projectHashPath = os.path.join(projectPath, ".pcc_project_id")
        self.projectHash = ""
        if not os.path.isfile(projectHashPath):
            self.projectHash =  base36.dumps(
                int(
                    hashlib.sha256(
                        self.HASH_SECRET + str(os.getuid()) + str(time.time())
                    ).hexdigest(),
                    16
                )
            )
            with open(projectHashPath, "w") as f:
                f.write(self.projectHash)
        if not self.projectHash:
            with open(projectHashPath, "r") as f:
                self.projectHash = f.read()
        self.vars = PlatformVars(self.projectHash)

    def getServices(self):
        """ Get list of service dependencies for project. """
        serviceConf = {}
        serviceList = []
        pathToServicesYaml = os.path.join(
            self.projectPath,
            PlatformServiceConfig.PLATFORM_SERVICES_PATH
        )
        with open(pathToServicesYaml, "r") as f:
            serviceConf = yaml.load(f, Loader=yamlordereddictloader.Loader)
        for serviceName in serviceConf:
            serviceList.append(
                PlatformService(
                    self.projectHash,
                    self.projectPath,
                    serviceName,
                    self.logger
                )
            )
        return serviceList

    def getApplications(self, withVars = True):
        """ Get all applications in project. """
        routerConfig = self.getRouterConfig()
        services = self.getServices()
        topPlatformAppConfigPath = os.path.join(self.projectPath, PlatformAppConfig.PLATFORM_FILENAME)
        projectVars = {}
        if withVars:
            projectVars = self.vars.all()
        if os.path.exists(topPlatformAppConfigPath):
            return [PlatformApp(self.projectHash, self.projectPath, services, projectVars, routerConfig, self.logger)]
        apps = []
        for path in os.listdir(os.path.realpath(self.projectPath)):
            path = os.path.join(self.projectPath, path)
            if os.path.isdir(path):
                platformAppConfigPath = os.path.join(path, PlatformAppConfig.PLATFORM_FILENAME)
                if os.path.isfile(platformAppConfigPath):
                    apps.append(PlatformApp(self.projectHash, self.projectPath, services, projectVars, routerConfig, self.logger))
        return apps

    def generateSshKey(self):
        """ Generate SSH key for use inside containers. """
        key = RSA.generate(2048)
        self.vars.set(
            "private_key",
            key.exportKey('PEM')
        )
        pubkey = key.publickey()
        self.vars.set(
            "public_key",
            pubkey.exportKey('OpenSSH')
        )

    def getProjectDomains(self):
        """ Get domains to use for this project. """
        projectDomains = self.vars.get(
            "project:domains"
        )
        projectDomains = projectDomains.strip().lower().split(",") if projectDomains else []
        projectDomains += [self.projectHash[:6]]
        return projectDomains

    def getRouterConfig(self, redirects = False):
        """ Parse routes.yaml config. """
        # open routes.yaml
        routeYamlPath = os.path.join(
            self.projectPath,
            self.PLATFORM_ROUTES_PATH
        )
        routes = {}
        if os.path.exists(routeYamlPath):
            with open(routeYamlPath, "r") as f:
                routes = yaml.load(f, Loader=yamlordereddictloader.Loader) 
        # generate config
        output = collections.OrderedDict()
        projectDomains = self.getProjectDomains()
        for routeSyntax in routes:
            parseRouteSyntax = urlparse(routeSyntax)
            isHttps = parseRouteSyntax.scheme == "https"
            serverKeys = [
                "%s://%s" % (
                    "https" if isHttps else "http",
                    parseRouteSyntax.hostname
                )
            ]
            # replace {DEFAULT} with project:domains
            for projectDomain in projectDomains:
                serverKey = serverKeys[0].replace(
                    self.ROUTE_DOMAIN_REPLACE,
                    projectDomain
                )
                if not serverKey in serverKeys:
                    serverKeys.append(serverKey)
            # generate internal domains for all server keys
            generatedServerKeys = []
            for serverKey in serverKeys:
                parseServerKey = urlparse(serverKey)
                generatedServerKeys.append(
                    "%s://%s%s%s" % (
                        "https" if parseServerKey.scheme == "https" else "http",
                        parseServerKey.hostname.replace(".", self.ROUTE_INTERNAL_DOMAIN_DOT_REPLACE),
                        self.ROUTE_INTERNAL_DOMAIN_DOT_REPLACE,
                        self.projectHash[:6]
                    )
                )
            for serverKey in serverKeys + generatedServerKeys:
                if not serverKey in output:
                    output[serverKey] = {
                        "type" :                routes[routeSyntax].get("type", "upstream"),
                        "upstream" :            routes[routeSyntax].get("upstream", "",).split(":")[0],
                        "to" :                  routes[routeSyntax].get("to", ""),
                        "cache" :               routes[routeSyntax].get("cache", {}),
                        "ssi" :                 routes[routeSyntax].get("ssi", {}),
                        "original_url" :        routeSyntax,
                        "redirects" :           routes[routeSyntax].get("redirects", {}) if redirects else {},
                        "is_platform_cc" :      True
                    }
        return output

    def generateRouterNginxConfig(self):
        """ Generate vhost config for nginx router. """
        apps = self.getApplications()
        projectDomains = self.getProjectDomains()
        nginxConf = ""
        for route, config in self.getRouterConfig(True).items():
            
            parseRoute = urlparse(route)
            isHttps = parseRoute.scheme == "https"

            nginxConf += "server {\n"
            # server name
            nginxConf += "\tserver_name %s;\n" % (
                parseRoute.hostname
            )
            # https
            if isHttps:
                nginxConf += "\tlisten 443 ssl;\n"
                nginxConf += "\tssl_certificate /etc/nginx/ssl/server.crt;\n"
                nginxConf += "\tssl_certificate_key /etc/nginx/ssl/server.key;\n"
            # http
            else:
                nginxConf += "\tlisten 80;\n"
            # main location
            nginxConf += "\tlocation %s {\n" % (
                "/%s" % parseRoute.path.lstrip("/")
            )
            # upstream
            if config.get("type", None) == "upstream":
                for app in apps:
                    if app.config.getName() == config.get("upstream", None):
                        # get app specific nginx config
                        nginxConf += str(app.docker.getProvisioner().generateNginxConfig())
                        break
                # redirects
                redirectPaths = config.get("redirects", {}).get("paths", {})
                if redirectPaths:
                    for location, redirectConfig in redirectPaths.items():
                        nginxConf += "\t}\n"
                        nginxConf += "\tlocation %s {\n" % (
                            location
                        )
                        nginxConf += "\t\treturn 301 %s$request_uri;\n" % (
                            redirectConfig.get("to", "/")
                        )

            # redirect
            elif config.get("type", None) == "redirect":
                nginxConf += "\t\treturn 301 %s$request_uri;\n" % (
                    config.get("to", "/").replace(
                        self.ROUTE_DOMAIN_REPLACE,
                        str(projectDomains[0])
                    )
                )
            # end server block
            nginxConf += "\t}\n"
            nginxConf += "}\n"

        return nginxConf

    def outputInfo(self, services = True, applications = True, routes = True):
        """ Output information about project. """
        if not self.logger: return

        self.logger.command.line(
            "\n======== Project '%s' =========\n" % self.projectHash[:6]
        )

        # display info about services
        if services:
            tableData = [
                ["Name", "Type", "Status", "IP Address"]
            ]
            for service in self.getServices():
                tableData.append([
                    service.config.getName(),
                    service.config.getType(),
                    service.docker.status(),
                    service.docker.getIpAddress() or "n/a"
                ])
            table = AsciiTable(tableData, "Services")
            self.logger.command.line(table.table)
            self.logger.command.line("")

        # display info about applications
        if applications:
            tableData = [
                ["Name", "Type", "Status", "IP Address"]
            ]
            for app in self.getApplications():
                tableData.append([
                    app.config.getName(),
                    app.config.getType(),
                    app.docker.status(),
                    app.docker.getIpAddress() or "n/a"
                ])
            table = AsciiTable(tableData, "Applications")
            self.logger.command.line(table.table)
            self.logger.command.line("")

        # display info about routes
        if routes:
            tableData = [
                ["Route", "Type", "Upstream / Redirect", "Original Url"]
            ]
            for route, routeConfig in self.getRouterConfig().items():
                to = routeConfig.get("to", "n/a")
                if routeConfig.get("type", "n/a") == "upstream":
                    to = routeConfig.get("upstream", "n/a")
                tableData.append([
                    route,
                    routeConfig.get("type", "n/a"),
                    to,
                    routeConfig.get("original_url", "")
                ])
            table = AsciiTable(tableData, "Routes")
            self.logger.command.line(table.table)
            self.logger.command.line("")

    def start(self):
        """ Start all services and apps. """
        for service in self.getServices():
            service.start()
        for app in self.getApplications():
            app.start()

    def stop(self):
        """ Stop all services and apps. """
        for app in self.getApplications():
            app.stop()
        for service in self.getServices():
            service.stop()

    def provision(self):
        """ Provision all services and apps. """
        if self.logger:
            self.logger.logEvent(
                "Provision services."
            )
        for service in self.getServices():
            service.docker.provision()
        for app in self.getApplications():
            app.provision()

    def deploy(self):
        """ Deploy all apps. """
        for app in self.getApplications():
            app.deploy()

    def purge(self):
        """ Purge all project data (including app volumes). """
        if self.logger:
            self.logger.logEvent(
                "Purge project '%s.' Waiting 5 seconds... (Press CTRL+C to cancel.)" % (
                     self.projectHash[:6]
                )
            )
        time.sleep(5)
        if self.logger:
            self.logger.logEvent("Purge start.")
        # itterate apps
        for app in self.getApplications():
            app.purge()
        # itterate services
        for service in self.getServices():
            service.stop()
            service.docker.purge()
        # purge vars
        if self.logger:
            self.logger.logEvent("Delete vars.")
        try:
            varsVolume = docker.from_env().volumes.get(self.vars.getVolumeKey())
            varsVolume.remove()
        except docker.errors.NotFound:
            pass

        