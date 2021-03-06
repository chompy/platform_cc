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

from __future__ import absolute_import
from future.moves.urllib.parse import urlparse
import os
import collections
import yaml
import yamlordereddictloader
from .base import BasePlatformParser
from ..exception.parser_error import ParserError


class RoutesParser(BasePlatformParser):
    """
    Routes (.platform/routes.yaml) parser.
    """

    """ Paths to routes yaml file. """
    YAML_PATHS = [
        ".platform/routes.yaml",
        ".platform/routes.pcc.yaml"
    ]

    def __init__(self, project, params={}):
        """
        :param project: Project data dictionary
        """
        BasePlatformParser.__init__(self, project.get("path"))
        self.project = project
        self._routes = None
        self.params = params

    def _readYaml(self, paths):
        config = {}
        if type(paths) is not list:
            paths = [paths]
        for path in paths:
            with open(path, "r") as f:
                self._mergeDict(
                    yaml.load(f, Loader=yamlordereddictloader.SafeLoader),
                    config
                )
        return config

    def getDefaultDomain(self):
        """
        Get default domain for current project.

        :return: Default domain
        :rtype: string
        """
        return self.project.get("short_uid", "default") + ".*"

    def getAllDomains(self):
        """
        Get all domains for current project.

        :return: List of domains
        :rtype: list
        """
        return [
            self.project.get("short_uid", "default") + ".*"
        ]

    def getRoutes(self):
        """
        Get list of all routes.

        :return: List of all route configurations
        :rtype: list
        """
        if self._routes:
            return self._routes
        yamlPaths = []
        for yamlPath in self.YAML_PATHS:
            yamlFullPath = os.path.join(
                self.projectPath,
                yamlPath
            )
            if os.path.isfile(yamlFullPath):
                yamlPaths.append(yamlFullPath)
        routes = self._readYaml(yamlPaths)
        self._routes = []
        for routeSyntax, routeConfig in routes.items():
            if routeSyntax.startswith("."):
                continue
            parseRouteSyntax = urlparse(routeSyntax)
            hostnames = []
            for domain in self.getAllDomains():
                newHostname = parseRouteSyntax.hostname.replace(
                    "{all}", domain
                )
                if newHostname != parseRouteSyntax.hostname:
                    hostnames.append(newHostname)
            newHostname = parseRouteSyntax.hostname.replace(
                "{default}", self.getDefaultDomain()
            )
            if newHostname != parseRouteSyntax.hostname:
                hostnames.append(newHostname)
            if not hostnames:
                if not self.params.get("disable_main_routes"):
                    hostnames.append(parseRouteSyntax.hostname)
                if not self.params.get("disable_extra_routes"):
                    extraSubDomainSep = self.params.get("extra_domain_seperator", ".")
                    hostnames.append(
                        "%s%s%s" % (
                            parseRouteSyntax.hostname.rstrip(".").replace(".", extraSubDomainSep),
                            extraSubDomainSep,
                            self.params.get("extra_domain_suffix", self.getDefaultDomain()).lstrip(".").lstrip("-")
                        )
                    )
            upstream = ["", ""]
            if routeConfig.get("type", "upstream") == "upstream":
                upstream = routeConfig.get("upstream", "").split(":")
                if len(upstream) < 2 or upstream[1][0:4] != "http":
                    continue
            self._routes.append({
                "scheme":             parseRouteSyntax.scheme,
                "hostnames":          hostnames,
                "path":               parseRouteSyntax.path,
                "type":               routeConfig.get("type", "upstream"),
                "upstream":           upstream[0],
                "to":                 routeConfig.get("to", ""),
                "cache":              routeConfig.get("cache", {}),
                "ssi":                routeConfig.get("ssi", {}),
                "original_url":       routeSyntax,
                "redirects":          routeConfig.get("redirects", {})
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
        output = collections.OrderedDict()

        # create upstreams to services
        if self.project and self.project.get("_enable_service_routes"):
            projectServices = self.project.get("services", {})
            for serviceName, serviceConf in projectServices.items():
                platformRelationships = serviceConf.get("platform_relationships", {})
                if not platformRelationships: continue
                firstPR = platformRelationships[list(platformRelationships.keys())[0]]
                if not firstPR.get("host", ""): continue
                path = "http://%s:%s" % (
                    firstPR.get("host", ""),
                    firstPR.get("port", "80")
                )
                output[path] = {
                    "type" : "upstream",
                    "upstream" : serviceName
                }

        for index, route in enumerate(routes):
            hostnames = route.pop("hostnames", [])
            route.pop("redirects")
            for hostname in hostnames:
                path = "%s://%s/%s" % (
                    route.get("scheme", "http"),
                    hostname,
                    route.get("path", "").lstrip("/")
                )
                output[path] = route
                output[path]["primary"] = (index == 0)
        return output
