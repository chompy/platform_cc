from __future__ import absolute_import
from future.moves.urllib.parse import urlparse
import os
import collections
import yaml
from .base import BasePlatformParser
from platform_cc.exception.parser_error import ParserError

class RoutesParser(BasePlatformParser):
    """
    Routes (.platform/routes.yaml) parser.
    """

    """ Path to routes yaml file. """
    YAML_PATH = ".platform/routes.yaml"

    def __init__(self, project):
        """
        :param project: Project data dictionary
        """
        BasePlatformParser.__init__(self, project.get("path"))
        self.project = project
        self._routes = None

    def _readYaml(self, path):
        config = None
        with open(path, "r") as f:
            config = yaml.load(f)
        return config

    def getDefaultDomain(self):
        """
        Get default domain for current project.

        :return: Default domain
        :rtype: string
        """
        return self.project.get("short_uid", "default")

    def getAllDomains(self):
        """
        Get all domains for current project.

        :return: List of domains
        :rtype: list
        """
        return [
            self.project.get("short_uid", "default")
        ]

    def getRoutes(self):
        """
        Get list of all routes.

        :return: List of all route configurations
        :rtype: list
        """
        if self._routes:
            return self._routes
        routes = self._readYaml(
            os.path.join(
                self.projectPath,
                self.YAML_PATH
            )
        )
        self._routes = []
        for routeSyntax, routeConfig in routes.items():
            if routeSyntax.startswith("."): continue
            parseRouteSyntax = urlparse(routeSyntax)
            hostnames = []
            for domain in self.getAllDomains():
                newHostname = parseRouteSyntax.hostname.replace("{all}", domain)
                if newHostname != parseRouteSyntax.hostname:
                    hostnames.append(newHostname)
            newHostname = parseRouteSyntax.hostname.replace("{default}", self.getDefaultDomain())
            if newHostname != parseRouteSyntax.hostname:
                hostnames.append(newHostname)
            if not hostnames:
                hostnames.append(parseRouteSyntax.hostname)
                hostnames.append(
                    "%s.%s.*" % (
                        parseRouteSyntax.hostname,
                        self.project.get("short_uid", "default")
                    )
                )
            upstream = ["", ""]
            if routeConfig.get("type", "upstream") == "upstream":
                upstream = routeConfig.get("upstream", "").split(":")
                if len(upstream) < 2 or upstream[1] != "http":
                    continue
            self._routes.append({
                "scheme"            : parseRouteSyntax.scheme,
                "hostnames"         : hostnames,
                "path"              : parseRouteSyntax.path,
                "type"              : routeConfig.get("type", "upstream"),
                "upstream"          : upstream[0],
                "to"                : routeConfig.get("to", ""),
                "cache"             : routeConfig.get("cache", {}),
                "ssi"               : routeConfig.get("ssi", {}),
                "original_url"      : routeSyntax,
                "redirects"         : routeConfig.get("redirects", {})
            })
        return self._routes

    def getRoutesByHostname(self):
        """
        Get all routes grouped by their hostname.

        :return: Dictionary of hostname to routes
        :rtype: dict
        """
        routes = self.getRoutes()
        output = {}
        for route in routes:
            for hostname in route.get("hostnames", []):
                if hostname not in output:
                    output[hostname] = []
                output[hostname].append(route)
        return output

    def getRoutesEnvironmentVariable(self):
        """
        Build PLATFORM_ROUTES environment variable data.

        :return: Dictionary of routes
        :rtype: dict
        """
        routes = self.getRoutes()
        output = {}        
        for route in routes:
            hostnames = route.pop("hostnames", [])
            route.pop("redirects")
            for hostname in hostnames:
                path = "%s://%s/%s" % (
                    route.get("scheme", "http"),
                    hostname,
                    route.get("path", "").lstrip("/")
                )
                output[path] = route
        return output
