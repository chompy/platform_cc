import io
import os
import json
from .base import BasePlatformApplication
from platform_cc.exception.container_command_error import ContainerCommandError

class PhpApplication(BasePlatformApplication):
    """
    Handler for PHP applications.
    """

    """ Mapping for application type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "php"             : "registry.gitlab.com/contextualcode/platform_cc/php56-fpm",
        "php:5.4"         : "registry.gitlab.com/contextualcode/platform_cc/php54-fpm",
        "php:5.6"         : "registry.gitlab.com/contextualcode/platform_cc/php56-fpm",
        "php:7.0"         : "registry.gitlab.com/contextualcode/platform_cc/php70-fpm",
        #"php:7.1"         : "registry.gitlab.com/contextualcode/platform_cc/php71-fpm",
        "php:7.2"         : "registry.gitlab.com/contextualcode/platform_cc/php72-fpm"   
    }

    """ Default user id to assign for user 'web' """
    DEFAULT_WEB_USER_ID = 1000

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
            sed -i "s/user = .*/user = web/g" /usr/local/etc/php-fpm.d/www.conf
            sed -i "s/group = .*/group = web/g" /usr/local/etc/php-fpm.d/www.conf
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
        # install extensions
        extInstall = self.config.get("runtime", {}).get("extensions", [])
        for extension in extInstall:
            if type(extension) is not str: continue
            self.logger.info(
                "Install/build '%s' extension.",
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
                    "web"
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
        def addFastCgi(scriptName = False):
            conf = ""
            if scriptName and type(scriptName) is str:
                conf += "\t\t\t\tset $_rewrite_path \"/%s\";\n" % scriptName.strip("/")
                conf += "\t\t\t\ttry_files $fastcgi_script_name @rewrite;\n"
            else:
                conf += "\t\t\t\ttry_files $fastcgi_script_name =404;\n"
            conf += "\t\t\t\tfastcgi_pass 127.0.0.1:9000;\n"
            conf += "\t\t\t\tset $_document_root $document_root;\n"
            conf += "\t\t\t\tinclude fastcgi_params;\n"
            conf += "\t\t\t\tfastcgi_split_path_info ^(.+?\.php)(/.*)$;\n"
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
            if type(passthru) is str:
                appNginxConf += "\t\t\tset $_rewrite_path \"/%s\";\n" % passthru.strip("/")
                appNginxConf += "\t\t\ttry_files $uri @rewrite;\n"
            else:
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
            if type(passthru) is str:
                appNginxConf += "\t\t\t\tset $_rewrite_path \"/%s\";\n" % passthru.strip("/")
                appNginxConf += "\t\t\t\ttry_files $uri @rewrite;\n"
            else:
                appNginxConf += "\t\t\t\ttry_files $uri =404;\n"
            appNginxConf += "\t\t\t\texpires -1s;\n"
            appNginxConf += "\t\t\t}\n"
            # == PASSTHRU
            passthru = locations[path].get("passthru", False)
            if passthru and not locations[path].get("scripts", False):
                if passthru == True: passthru = "/index.php"
                appNginxConf += "\t\t\tlocation ~ \".+?\.php(?=$|/)\" {\n"
                appNginxConf += "\t\t\t\tallow all;\n"
                appNginxConf += addFastCgi(passthru)
                appNginxConf += "\t\t\t}\n"
                #appNginxConf += "\t\tlocation / {\n"
                #appNginxConf += "\t\t\ttry_files $uri /%s$is_args$args;\n" % passthru.strip("/")
                #appNginxConf += "\t\t}\n"
            # == SCRIPTS
            scripts = locations[path].get("scripts", False)
            if scripts:
                appNginxConf += "\t\t\tlocation ~ [^/]\.php(/|$) {\n"
                appNginxConf += addFastCgi()
                if passthru:
                    appNginxConf += "\t\t\t\tfastcgi_index %s;\n" % (passthru.lstrip("/"))
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
            """
            rules = locations[path].get("rules", {})
            if rules:
                for ruleRegex in rules:
                    rule = rules[ruleRegex]
                    appNginxConf += "\t\t\tlocation ~ %s {\n" % (ruleRegex)
                    # allow
                    if not rule.get("allow", True):
                        appNginxConf += "\t\t\t\tdeny all;\n"
                    else:
                        appNginxConf += "\t\t\t\tallow all;\n"
                    # passthru
                    passthru = rule.get("passthru", False)
                    if passthru:
                        appNginxConf += addFastCgi(passthru)
                    # expires
                    expires = rule.get("expires", False)
                    if expires:
                        appNginxConf += "\t\t\t\texpires %s;\n" % expires
                    # headers
                    headers = rule.get("headers", {})
                    for headerName in headers:
                        appNginxConf += "\t\t\t\tadd_header %s %s;\n" % (
                            headerName,
                            headers[headerName]
                        )
                    # scripts
                    scripts = rule.get("scripts", False)
                    appNginxConf += "\t\t\t\tlocation ~ [^/]\.php(/|$) {\n"
                    if scripts:
                        appNginxConf += addFastCgi()
                        if passthru:
                            appNginxConf += "\t\t\t\t\tfastcgi_index %s;\n" % (passthru.lstrip("/"))
                    else:
                        appNginxConf += "\t\t\t\t\tdeny all;\n"
                    appNginxConf += "\t\t\t\t}\n"
            """
            appNginxConf += "\t\t}\n"
        return appNginxConf

    def start(self, requireServices = True):
        BasePlatformApplication.start(self, requireServices)
        container = self.getContainer()
        if not container: return
        # link php.ini in app root
        self.logger.info(
            "Compile php.ini."
        )
        self.runCommand(
            "[ -f /app/php.ini ] && ln -s -f /app/php.ini /usr/local/etc/php/conf.d/02-app.ini || true"
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
