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
from nginx.config.api import Location
from nginx.config.api.options import KeyValueOption, KeyValuesMultiLines, KeyOption
from .base import BasePlatformApplication
from platform_cc.exception.container_command_error import ContainerCommandError
from platform_cc.version import PCC_VERSION

class PhpApplication(BasePlatformApplication):
    """
    Handler for PHP applications.
    """

    """ Mapping for application type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "php"             : "chompy/platform_cc:%s-php73" % PCC_VERSION,
        "php:5.4"         : "chompy/platform_cc:%s-php54" % PCC_VERSION,
        "php:5.6"         : "chompy/platform_cc:%s-php56" % PCC_VERSION,
        "php:7.0"         : "chompy/platform_cc:%s-php70" % PCC_VERSION,
        "php:7.1"         : "chompy/platform_cc:%s-php71" % PCC_VERSION,
        "php:7.2"         : "chompy/platform_cc:%s-php72" % PCC_VERSION,
        "php:7.3"         : "chompy/platform_cc:%s-php73" % PCC_VERSION,
    }

    """ Default user id to assign for user 'web' """
    DEFAULT_WEB_USER_ID = 1000

    TCP_PORT = 9000

    """ Path to PHP extension configuration JSON. """
    EXTENSION_CONF_JSON = os.path.join(
        os.path.dirname(__file__),
        "../data/php_extensions.json"
    )

    """ Path to PHP nginx configuration. """
    NGINX_CONF = os.path.join(
        os.path.dirname(__file__),
        "../data/php_nginx.conf"
    )

    def getContainerCommand(self):
        if self.getDockerImage() == self.getBaseImage():
            return None
        command = self.config.get("web", {}).get("commands", {}).get("start")
        if command:
            return "sh -c \"%s\"" % command
        return None

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getExtensionInstallCommand(self, extensionName):
        """
        Get command needed to install a given extension.

        :param extensionName: Name of extension
        :return: Command to run
        :rtype: str
        """
        try:
            self._extensionConf
        except AttributeError:
            self._extensionConf = None
        if not self._extensionConf:
            if not os.path.exists(self.EXTENSION_CONF_JSON):
                return ""
            with open(self.EXTENSION_CONF_JSON, "r") as f:
                self._extensionConf = json.load(f)
        extensionConfEntry = self._extensionConf.get(
            extensionName,
            self._extensionConf.get("__default__", None)
        )
        if not extensionConfEntry: return ""
        for imageConf in extensionConfEntry:
            if self.getBaseImage() not in imageConf.get("images", []):
                continue
            command = imageConf.get("command", "").replace("__EXT_NAME__", extensionName)
            return command
        return ""

    def build(self):
        self.prebuild()
        output = ""
        # change web user id
        userId = self.project.get("config", {}).get("web_user_id", self.DEFAULT_WEB_USER_ID)
        if userId != self.DEFAULT_WEB_USER_ID:
            self.logger.info(
                "Update 'web' user id."
            )        
            output += self.runCommand(
                """
                usermod -u %s web
                """ % (
                    userId
                )
            )
        # install ssh key + known_hosts
        self.installSsh()
        output += self.runCommand(
            "chown -f -R web /app/.ssh"
        )
        # install extensions
        extInstall = self.config.get("runtime", {}).get("extensions", [])
        output += self.runCommand(
            """
            apt-get update -y
            """
        )
        for extension in extInstall:
            if type(extension) is not str: continue
            self.logger.info(
                "Enable '%s' extension.",
                extension
            )
            command = self.getExtensionInstallCommand(extension)
            if not command: continue
            output += self.runCommand(command)
        # build flavor composer
        if self.config.get("build", {}).get("flavor") == "composer":
            self.logger.info(
                "Composer install."
            )
            try:
                output += self.runCommand(
                    """
                    php -d memory_limit=-1 /usr/local/bin/composer install
                    """,
                    "root"
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
                "root"
            )
        # allow build hooks to fail...for now
        except ContainerCommandError:
            pass
        # attempt to fix file permissions
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
        self.stop()
        return output

    def _generateNginxPassthruOptions(self, locationConfig = {}, script = False):
        
        # force fastcgi/tcp upstream for php
        if not "web" in self.config:
            self.config["web"] = {}
        if not "upstream" in self.config["web"]:
            self.config["web"]["upstream"] = {}
        if not "socket_family" in self.config["web"]["upstream"]:
            self.config["web"]["upstream"]["socket_family"] = "socket"
        if not "protocol" in self.config["web"]["upstream"]:
            self.config["web"]["upstream"]["protocol"] = "fastcgi"            

        options = BasePlatformApplication._generateNginxPassthruOptions(self, locationConfig)
        setOptions = [
            "$_document_root $document_root",
            "$path_info $fastcgi_path_info"
        ]
        if script:
            script = str(script)
            setOptions.append(
                "$_rewrite_path \"/%s\"" % script.strip("/")
            )
            options.append(
                KeyValueOption("try_files", "$fastcgi_script_name @rewrite")
            )
        else:
            options.append(
                KeyValueOption("try_files", "$fastcgi_script_name =404")
            )
        options.append(
            KeyValuesMultiLines("set", setOptions)
        )
        options.append(
            KeyValueOption("fastcgi_split_path_info", "^(.+?\.php)(/.*)$")
        )
        return options

    def _generateNginxLocations(self, path, locationConfig = {}):

        # params
        pathStrip = "/%s/" % path.strip("/")
        if pathStrip == "//": pathStrip = "/"
        passthru = locationConfig.get("passthru", False)
        if passthru == True: passthru = "/index.php"
        if passthru != False: passthru = str(passthru)
        scripts = locationConfig.get("scripts", False)
        index = locationConfig.get("index", [])
        if type(index) is not list: index = [index]

        # get base locations
        locations = BasePlatformApplication._generateNginxLocations(self, path, locationConfig)

        # update root location
        if passthru:
            locations[0].options["try_files"] = "$uri @rewrite"
        locations[0].options["set"] = ("$_rewrite_path \"/%s\"" % passthru.strip("/")) if passthru else "$_rewrite_path \"\""

        # update main location
        # php specific passthru
        if passthru:
            locations[1].sections.pop("location ~ /")
            locations[1].sections.add(
                Location(
                    "~ \".+?\.php(?=$|/)\"",
                    allow = "all",
                    *self._generateNginxPassthruOptions(locationConfig, passthru)
                )
            )

        # php sub location
        subLocationOptions = {}
        if index:
            subLocationOptions["index"] = " ".join(index)
        if passthru:
            subLocationOptions["set"] = "$_rewrite_path \"/%s\"" % passthru.strip("/")
            subLocationOptions["try_files"] = "$uri @rewrite"
        locations[1].sections.add(
            Location(
                pathStrip,
                **subLocationOptions
            )
        )

        # php scripts
        if scripts:
            options = self._generateNginxPassthruOptions(locationConfig)
            if passthru:
                options.append(
                    KeyValueOption("fastcgi_index", passthru.lstrip("/"))
                )
            locations.append(
                Location(
                    "~ [^/]\\.php(/|$)",
                    *options
                )
            )

        return locations

    def startServices(self):
        BasePlatformApplication.startServices(self)
        # start newrelic agent
        extInstall = self.config.get("runtime", {}).get("extensions", [])
        if "newrelic" in extInstall:
            self.logger.info(
                "Start Newrelic agent."
            )
            self.runCommand(
                """
                newrelic-daemon
                """
            ) 

    def start(self, requireServices = True):
        BasePlatformApplication.start(self, requireServices)
        container = self.getContainer()
        if not container: return
        # link php.ini in app root
        self.logger.info(
            "Compile php.ini."
        )
        self.runCommand(
            "[ -f /app/php.ini ] && ln -s -f /app/php.ini /usr/local/etc/php/conf.d/zzz-03-app.ini || true"
        )
        # build php.ini from config vars
        phpIniConfig = self.config.get("variables", {}).get("php", {})
        phpIniFileObj = io.BytesIO()
        for key, value in phpIniConfig.items():
            phpIniFileObj.write(
                bytes(str("%s = %s\n" % (key, value)).encode("utf-8"))
            )
        for key, value in self.project.get("variables", {}).items():
            if not key or key[0:4] != "php:": continue
            phpIniFileObj.write(
                bytes(str("%s = %s\n" % (key[4:], value)).encode("utf-8"))
            )
        self.uploadFile(
            phpIniFileObj,
            "/usr/local/etc/php/conf.d/03-app.ini"
        )
        # restart container to reload conf changes
        container.restart()
        self.setupMounts()
        # start container services
        self.startServices()