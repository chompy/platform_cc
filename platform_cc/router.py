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

import os
import docker
import logging
from platform_cc.container import Container
from platform_cc.parser.routes import RoutesParser
from nginx.config.api import Location, Block
from nginx.config.api.options import KeyValueOption, KeyValuesMultiLines, KeyOption

class PlatformRouter(Container):

    """
    Main router for accessing all projects via the web.
    """

    """ Path to main nginx configuration. """
    NGINX_CONF = os.path.join(
        os.path.dirname(__file__),
        "data/router_nginx.conf"
    )

    """ Path to Nginx project config files inside container. """
    NGINX_PROJECT_CONF_PATH = "/router"

    def __init__(self, dockerClient = None):
        Container.__init__(self, {}, "router", dockerClient)
        self.logger = logging.getLogger(__name__)

    def getBaseImage(self):
        return "nginx:1.13"

    def getCommitImage(self):
        return "%s:%s" % (
            self.COMMIT_REPOSITORY_NAME,
            self.getName()
        )

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

    def generateNginxConfig(self, applications):
        """
        Generate Nginx vhost for applications in a project.

        :param applications: List of all applications in a project
        :return: Nginx configuration
        :rtype: str
        """
        self.logger.info(
            "Generate router Nginx configuration for project '%s.'.",
            applications[0].project.get("short_uid")
        )
        routesParser = RoutesParser(applications[0].project)
        routeHostnames = routesParser.getRoutesByHostname()
        output = ""
        for hostname, routes in routeHostnames.items():
            self.logger.info(
                "Add %s route(s) for '%s.'",
                len(routes),
                hostname
            )
            # create vhost entry for each scheme
            for scheme in ["http", "https"]:

                # create server section
                server = Block(
                    "server",
                    resolver = "127.0.0.11",
                    server_name = hostname,
                    listen = "443 ssl" if scheme == "https" else "80",
                    client_max_body_size = "200M"
                )

                # add ssl
                if scheme == "https":
                    server.options["ssl_certificate"] = "/etc/nginx/ssl/server.crt"
                    server.options["ssl_certificate_key"] = "/etc/nginx/ssl/server.key"

                # add locations
                hasRouteForScheme = False
                for config in routes:
                    if config.get("scheme", "http") != scheme: continue
                    hasRouteForScheme = True
                    path = config.get("path", "/")
                    if not path: path = "/"
                    location = Location(
                        path.replace(" ", '[\s]')
                    )
                    # type 'upstream'
                    if config.get("type") == "upstream":
                        # redirects
                        redirectHasRootPath = False
                        redirectPaths = config.get("redirects", {}).get("paths", {})
                        for _location, redirectConfig in redirectPaths.items():
                            if _location.strip() == "/":
                                redirectHasRootPath = True
                            location.sections.add(
                                Location(
                                    "~ %s" % _location.replace(" ", '[\s]'),
                                    KeyValueOption("return", "30 %s" % redirectConfig.get("to", "/"))
                                )
                            )
                        # upstream, proxy_pass
                        upstreamHost = ""
                        for application in applications:
                            if application.getName() == config.get("upstream"):
                                upstreamHost = application.getContainerName() # container host name
                        if not redirectHasRootPath and upstreamHost:
                            location.sections.add(
                                Location(
                                    "~* /",
                                    KeyValuesMultiLines(
                                        "proxy_set_header",
                                        [
                                            "X-Client-IP $server_addr",
                                            "X-Forwarded-Host $host",
                                            "X-Forwarded-Port $server_port",
                                            "X-Forwarded-Proto $scheme",
                                            "X-Forwarded-Server $host",
                                            "Host $host",
                                            "X-Forwarded-For $remote_addr",
                                        ]
                                    ),
                                    set = "$upstream http://%s" % upstreamHost,
                                    proxy_pass = "$upstream"
                                )
                            )
                    
                    # type 'redirect'
                    elif config.get("type") == "redirect":
                        to = config.get("to", "").replace("{default}", routesParser.getDefaultDomain())
                        location.sections.add(
                            Location(
                                "~ /",
                                KeyValueOption("return", ("501" if "*" in to else ("301 %s" % to)))
                            )
                        )

                    # add location to server
                    server.sections.add(location)                        

                # if scheme does not have any routes create a redirect
                if not hasRouteForScheme:
                    server.sections.add(
                        Location(
                            "/",
                            KeyValueOption(
                                "return",
                                "301 %s://$host$request_uri" % ("http" if scheme == "https" else "https")
                            )
                        )
                    )
                
                # add server output
                output += str(server)
        
        return output
    
    def build(self):
        # create web user and install dev ssl certificate
        self.logger.info(
            "Create 'web' user and create default SSL certificate."
        )
        self.runCommand(
            """
            mkdir %s
            useradd -d  /router -m -p secret~ web
            apt-get update
            apt-get install openssl
            apt-get clean
            mkdir /etc/nginx/ssl
            cd /etc/nginx/ssl
            openssl genrsa -des3 -passout "pass:^nx/{Dm[[k3b]ATf" -out server.pass.key 2048
            openssl rsa -passin "pass:^nx/{Dm[[k3b]ATf" -in server.pass.key -out server.key
            rm server.pass.key
            openssl req -new -key server.key -out server.csr \
                -subj "/C=US/ST=Florida/L=Tallahassee/O=ContextualCode/OU=Developers/CN=dev.local"
            openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
            """ % (self.NGINX_PROJECT_CONF_PATH)
        )
        # add router nginx.conf
        self.logger.info(
            "Add main Nginx configuration."
        )
        if os.path.exists(self.NGINX_CONF):
            with open(self.NGINX_CONF, "rb") as f:
                self.uploadFile(
                    f,
                    "/etc/nginx/nginx.conf"
                )
        # commit container
        self.commit()
        
    def start(self):
        Container.start(self)
        if self.getDockerImage() == self.getBaseImage():
            self.build()
            self.stop()
            return self.start()

    def restart(self):
        # restart without deleting the container
        container = self.getContainer()
        if not container:
            return self.start()
        self.logger.info("Restart router.")
        container.restart()
