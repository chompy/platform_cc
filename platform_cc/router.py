import os
import docker
import logging
from platform_cc.container import Container
from platform_cc.parser.routes import RoutesParser

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
                # start vhost
                output += "server {\n"
                # resolver
                output += "\tresolver 127.0.0.11;\n"
                # server_name
                output += "\tserver_name %s;\n" % (
                    hostname
                )
                # listen
                if scheme == "https":
                    output += "\tlisten 443 ssl;\n"
                    output += "\tssl_certificate /etc/nginx/ssl/server.crt;\n"
                    output += "\tssl_certificate_key /etc/nginx/ssl/server.key;\n"
                else:
                    output += "\tlisten 80;\n"
                # client_max_body_size
                output += "\tclient_max_body_size 200M;\n"
                # add locations
                hasRouteForScheme = False
                for config in routes:
                    if config.get("scheme", "http") != scheme: continue
                    hasRouteForScheme = True
                    # location
                    path = config.get("path", "/")
                    if not path: path = "/"
                    output += "\tlocation %s {\n" % (
                        path
                    )
                    # type 'upstream'
                    if config.get("type") == "upstream":
                        # redirects
                        redirectHasRootPath = False
                        redirectPaths = config.get("redirects", {}).get("paths", {})
                        for location, redirectConfig in redirectPaths.items():
                            if location.strip() == "/":
                                redirectHasRootPath = True
                            output += "\t\tlocation ~ %s {\n" % (
                                location
                            )
                            output += "\t\t\treturn 301 %s$request_uri;\n" % (
                                redirectConfig.get("to", "/")
                            )
                            output += "\t\t}\n"
                        # upstream, proxy_pass
                        upstreamHost = ""
                        for application in applications:
                            if application.getName() == config.get("upstream"):
                                upstreamHost = application.getContainerName() # container host name
                        if not redirectHasRootPath and upstreamHost:
                            output += "\t\tlocation ~* / {\n"
                            output += "\t\t\tset $upstream http://%s;\n" % (
                                upstreamHost
                            )
                            output += "\t\t\tproxy_set_header X-Forwarded-Host $host:$server_port;\n"
                            output += "\t\t\tproxy_set_header X-Forwarded-Proto $scheme;\n"
                            output += "\t\t\tproxy_set_header X-Forwarded-Server $host;\n"
                            output += "\t\t\tproxy_set_header X-Forwarded-For $remote_addr;\n"
                            output += "\t\t\tproxy_pass $upstream;\n"
                            output += "\t\t}\n"
                    # type 'redirect'
                    elif config.get("type") == "redirect":
                        output += "\t\tlocation ~ / {\n"
                        output += "\t\t\treturn 301 %s$request_uri;\n" % (
                            config.get("to", "").replace("{default}", routesParser.getDefaultDomain())
                        )
                        output += "\t\t}\n"
                    output += "\t}\n"
                # if scheme does not have any routes create a redirect
                if not hasRouteForScheme:
                    output += "\tlocation / {\n"
                    output += "\t\treturn 301 %s://$host$request_uri;\n" % (
                        ("http" if scheme == "https" else "https")
                    )
                    output += "\t}\n"
                output += "}\n"
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