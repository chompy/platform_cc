import os
import time
import hashlib
import base36
import yaml
import yamlordereddictloader
import collections
import docker
from urlparse import urlparse
from Crypto.PublicKey import RSA
from terminaltables import AsciiTable
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

    def getApplications(self, withVars = True):
        """ Get all applications in project. """
        topPlatformAppConfigPath = os.path.join(self.projectPath, PlatformAppConfig.PLATFORM_FILENAME)
        projectVars = {}
        if withVars:
            projectVars = self.vars.all()
        if os.path.exists(topPlatformAppConfigPath):
            return [PlatformApp(self.projectHash, self.projectPath, projectVars, self.logger)]
        apps = []
        for path in os.listdir(os.path.realpath(self.projectPath)):
            path = os.path.join(self.projectPath, path)
            if os.path.isdir(path):
                platformAppConfigPath = os.path.join(path, PlatformAppConfig.PLATFORM_FILENAME)
                if os.path.isfile(platformAppConfigPath):
                    apps.append(PlatformApp(self.projectHash, self.projectPath, projectVars, self.logger))
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

    def generateRouterConfig(self):
        """ Generate vhost config for nginx router. """

        routeYamlPath = os.path.join(
            self.projectPath,
            self.PLATFORM_ROUTES_PATH
        )
        routes = {}
        if os.path.exists(routeYamlPath):
            with open(routeYamlPath, "r") as f:
                routes = yaml.load(f, Loader=yamlordereddictloader.Loader)
        projectDomains = self.vars.get(
            "project:domains"
        )
        projectDomains = projectDomains.strip().lower().split(",") if projectDomains else []
        projectDomains += ["%s.local" % self.projectHash[:6]]
        serverList = collections.OrderedDict()
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
            nginxConf += "\tserver {\n"
            hostnames = []
            for projectDomain in projectDomains:
                hostname = serverList[serverName]["hostname"].replace(
                    self.ROUTE_DOMAIN_REPLACE,
                    projectDomain
                )
                if hostname not in hostnames:
                    hostnames.append(hostname)
            nginxConf += "\t\tserver_name %s;\n" % (
                str(" ".join(hostnames))
            )
            
            # TODO HTTPS
            nginxConf += "\t\tlisten 80;\n"

            paths = serverList[serverName]["paths"]
            for path in paths:
                nginxConf += "\t\tlocation %s {\n" % (
                    path
                )
                if paths[path]["type"] == "upstream":
                    upstream = paths[path]["upstream"].split(":")[0]
                    for app in self.getApplications():
                        if app.config.getName() == upstream:
                            import json
                            ipAddress = app.web.docker.getIpAddress()
                            nginxConf += "\t\t\tproxy_set_header X-Forwarded-Host $host:$server_port;\n"
                            nginxConf += "\t\t\tproxy_set_header X-Forwarded-Server $host;\n"
                            nginxConf += "\t\t\tproxy_set_header X-Forwarded-For $remote_addr;\n"
                            nginxConf += "\t\t\tproxy_pass http://%s;\n" % (
                                ipAddress
                            )
                            break
                elif paths[path]["type"] == "redirect":
                    to = paths[path].get("to", None)
                    if to:
                        nginxConf += "\t\t\treturn 301 %s$request_uri;\n" % (
                            to.replace(
                                self.ROUTE_DOMAIN_REPLACE,
                                str(projectDomains[0])
                            )
                        )
                nginxConf += "\t\t}\n"
            nginxConf += "\t}\n"

        return nginxConf

    def outputInfo(self):
        """ Output information about project. """
        if not self.logger: return

        tableData = [
            ["Application Name", "Status", "IP Address (Web)", "Services"]
        ]
        for app in self.getApplications():

            serviceNames = ""
            for service in app.getServices():
                serviceNames += "%s, " % (service.config.getName())

            tableData.append([
                app.config.getName(),
                app.docker.status(),
                app.web.docker.getIpAddress() or "n/a",
                serviceNames.strip().rstrip(",")
            ])
        table = AsciiTable(tableData, "Project '%s'" % self.projectHash[:6])
        self.logger.command.line(table.table)