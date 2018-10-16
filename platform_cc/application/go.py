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

import io
import os
import json
from .base import BasePlatformApplication
from platform_cc.exception.container_command_error import ContainerCommandError

class GoApplication(BasePlatformApplication):
    """
    Handler for Golang applications.
    """

    """ Mapping for application type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "golang:1.11"         : "registry.gitlab.com/contextualcode/platform_cc/golang-1-11",  
    }

    """ Default user id to assign for user 'web' """
    DEFAULT_WEB_USER_ID = 1000

    """ Port to use for TCP upstream. """
    TCP_PORT = 8001

    """ Socket path to use for upstream. """
    SOCKET_PATH = "/tmp/app.socket"

    def getContainerCommand(self):
        if self.getDockerImage() == self.getBaseImage():
            return None
        command = self.config.get("web", {}).get("commands", {}).get("start")
        if command:
            return "sh -c \"%s\"" % command
        return None

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerEnvironmentVariables(self):
        envVars = BasePlatformApplication.getContainerEnvironmentVariables(self)
        envVars["PORT"] = self.TCP_PORT
        envVars["SOCKET"] = self.SOCKET_PATH
        return envVars

    def build(self):
        self.prebuild()
        output = ""
        # add web user
        self.logger.info(
            "Add and configure 'web' user."
        )
        output += self.runCommand(
            """
            useradd -d /app -m -p secret~ --uid %s web
            usermod -a -G staff web
            mkdir -p /var/lib/gems
            chown -R web:web /var/lib/gems
            chown -R root:staff /usr/bin
            chmod -R g+rw /usr/bin
            """ % (
                self.project.get("config", {}).get("web_user_id", self.DEFAULT_WEB_USER_ID)
            )
        )
        output += self.runCommand(
            "usermod -u %s web" % (
                self.project.get("config", {}).get("web_user_id", self.DEFAULT_WEB_USER_ID)
            )            
        )
        # install ssh key + known_hosts
        self.installSsh()
        output += self.runCommand(
            "chown -f -R web /app/.ssh"
        )
        self.logger.info(
            "Setup/fix user permission."
        )
        try:
            output += self.runCommand(
                """
                chown -f -R web %s
                chown -f -R web %s
                """ % (self.STORAGE_DIRECTORY, self.APPLICATION_DIRECTORY)
            )
        except ContainerCommandError:
            pass
        # build hooks
        self.logger.info(
            "Run build hooks."
        )
        try:
            output += self.runCommand(
                self.config.get("hooks", {}).get("build", ""),
                "web"
            )
        # allow build hooks to fail...for now
        except ContainerCommandError:
            pass
        # clean up
        self.logger.info(
            "Clean up."
        )
        output += self.runCommand(
            """
            apt-get clean
            """
        )
        # commit container
        self.logger.info(
            "Commit container."
        )
        self.commit()
        return output

    def generateNginxConfig(self):
        """
        Generate configuration for nginx specific to application.

        :return: Nginx configuration
        :rtype: str
        """
        self.logger.info(
            "Generate application Nginx configuration."
        )
        locations = self.config.get("web", {}).get("locations", {})
        appNginxConf = ""
        def addPassthru():
            conf = ""
            upstreamConf = self.config.get("web", {}).get("upstream", {"socket_family" : "tcp", "protocol" : "http"})
            # tcp port, proxy pass
            if upstreamConf.get("socket_family") == "tcp" and upstreamConf.get("protocol") == "http":
                conf += "\t\t\t\tproxy_pass http://127.0.0.1:%d\n;" % self.TCP_PORT
                conf += "\t\t\t\tproxy_set_header Host $host;\n"
            # tcp port, fastcgi
            elif upstreamConf.get("socket_family") == "tcp" and upstreamConf.get("protocol") == "fastcgi":
                conf += "\t\t\t\tfastcgi_pass 127.0.0.1:%d;" % self.TCP_PORT
                conf += "\t\t\t\tinclude fastcgi_params;\n"
                conf += "\t\t\t\tset $path_info  $fastcgi_path_info;\n";
            # socket, proxy pass
            elif upstreamConf.get("socket_family") == "socket" and upstreamConf.get("protocol") == "http":
                conf += "\t\t\t\tproxy_pass unix:%s;\n" % self.SOCKET_PATH
                conf += "\t\t\t\tproxy_set_header Host $host;\n"
            # socket, fastcgi
            elif upstreamConf.get("socket_family") == "socket" and upstreamConf.get("protocol") == "fastcgi":
                conf += "\t\t\t\tfastcgi_pass unix:%s;" % self.SOCKET_PATH
                conf += "\t\t\t\tinclude fastcgi_params;\n"
                conf += "\t\t\t\tset $path_info  $fastcgi_path_info;\n";
            return conf

        for path in locations:
            root = locations[path].get("root", "") or ""
            passthru = locations[path].get("passthru", False)
            # ============
            appNginxConf += "\t\tlocation = \"%s\" {\n" % path.rstrip("/")
            appNginxConf += "\t\t\talias \"%s\";\n" % (
                ("%s/%s" % (self.APPLICATION_DIRECTORY, root.strip("/"))).rstrip("/")
            )
            appNginxConf += "\t\t\ttry_files $uri =404;\n"
            appNginxConf += "\t\t\texpires -1s;\n"
            appNginxConf += "\t\t}\n"
            # ============
            pathStrip = "/%s/" % path.strip("/")
            if pathStrip == "//": pathStrip = "/"
            appNginxConf += "\t\tlocation \"%s\" {\n" % pathStrip
            # == ALIAS
            appNginxConf += "\t\t\talias \"%s/\";\n" % (
                ("%s/%s" % (self.APPLICATION_DIRECTORY, root.strip("/"))).rstrip("/")
            )
            # == HEADERS
            headers = locations[path].get("headers", {})
            for headerName in headers:
                appNginxConf += "\t\t\tadd_header %s %s;\n" % (
                    headerName,
                    headers[headerName]
                )
            # == SUB LOCATION
            appNginxConf += "\t\t\tlocation \"%s\" {\n" % pathStrip
            appNginxConf += "\t\t\t\texpires -1s;\n"
            appNginxConf += "\t\t\t}\n"
            # == PASSTHRU
            passthru = locations[path].get("passthru", False)
            if passthru:
                appNginxConf += "\t\t\tlocation ~ / {\n"
                appNginxConf += "\t\t\t\tallow all;\n"
                appNginxConf += addPassthru()
                appNginxConf += "\t\t\t}\n"
            # == ALLOW
            #allow = locations[path].get("allow", False)
            # TODO!
            # allow = false should deny access when requesting a file that does exist but
            # does not match a rule
            # == RULES
            # TODO
            # we don't currently make use of the rules directive, so this code has
            # not been tested, commented out for now
            # see php
            appNginxConf += "\t\t}\n"
        return appNginxConf

    def start(self, requireServices = True):
        BasePlatformApplication.start(self, requireServices)
        container = self.getContainer()
        if not container: return
        # setup mount points
        self.setupMounts()
        # not yet built/provisioned
        if self.getDockerImage() == self.getBaseImage():
            self.build()
            return self.start(requireServices)
        # nginx config
        nginxConfFileObj = io.BytesIO(
            bytes(str(self.generateNginxConfig()).encode("utf-8"))
        )
        self.uploadFile(
            nginxConfFileObj,
            "/usr/local/nginx/conf/app.conf"
        )
        # start nginx + other services
        self.logger.info(
            "Start Nginx."
        )
        self.runCommand(
            """
            /usr/local/nginx/sbin/nginx
            """
        )

        # install cron jobs if enabled
        self.installCron()
