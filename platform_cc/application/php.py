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
        "php"             : "php:5.6-fpm",
        "php:5.4"         : "php:5.4-fpm",
        "php:5.6"         : "php:5.6-fpm",
        "php:7.0"         : "php:7.0-fpm"   
    }

    """ Default UID to assign for user 'web' """
    DEFAULT_WEB_UID = 1000

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
        output = ""
        # create 'web' user
        self.logger.info(
            "Create 'web' user."
        )
        output += self.runCommand(
            "useradd -d /app -m -p secret~ --uid %s web || true" % (
                self.project.get("config", {}).get("web_uid", self.DEFAULT_WEB_UID)
            )            
        )
        # provision container
        self.logger.info(
            "Install/build dependencies."
        )
        output += self.runCommand(
            """
            usermod -a -G staff web
            chown -R web:web /app
            apt-get update
            apt-get install -y rsync git unzip python-pip python-dev \\
                gem nodejs npm libyaml-dev ruby ruby-dev nginx less nano \\
                libicu-dev libxslt1-dev libfreetype6-dev libjpeg62-turbo-dev libpng12-dev \\
                libmcrypt-dev
            mkdir -p /var/lib/gems
            chown -R web:web /var/lib/gems
            chown -R root:staff /usr/bin
            chmod -R g+rw /usr/bin
            ln -s /usr/bin/nodejs /usr/bin/node
            sed -i "s/user = .*/user = web/g" /usr/local/etc/php-fpm.d/www.conf
            sed -i "s/group = .*/group = web/g" /usr/local/etc/php-fpm.d/www.conf
            echo "date.timezone = UTC" > /usr/local/etc/php/conf.d/01-main.ini
            echo "memory_limit = 512M" >> /usr/local/etc/php/conf.d/01-main.ini
            ln -s /usr/local/sbin/php-fpm /usr/sbin/php5-fpm
            ln -s /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.0
            ln -s /usr/local/sbin/php-fpm /usr/sbin/php-fpm7.1-zts
            php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
            php composer-setup.php --install-dir=/usr/local/bin
            rm composer-setup.php
            ln -s /usr/local/bin/composer.phar /usr/local/bin/composer
            docker-php-ext-install -j$(nproc) bcmath intl xsl mysql mysqli pdo_mysql sockets exif mcrypt
            docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/
            docker-php-ext-install -j$(nproc) gd
            chown -R web:web %s
            chown -R web:web %s
            """ % (self.STORAGE_DIRECTORY, self.APPLICATION_DIRECTORY)
        )
        # install nginx config
        self.logger.info(
            "Install main Nginx configuration file."
        )
        with open(self.NGINX_CONF, "rb") as f:
            self.uploadFile(f, "/etc/nginx/nginx.conf")
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
        def addFastCgi(scriptName = ""):
            if not scriptName: scriptName = "$fastcgi_script_name"
            conf = ""
            conf += "\t\t\t\tfastcgi_split_path_info ^(.+?\.php)(/.*)$;\n"
            conf += "\t\t\t\tfastcgi_pass 127.0.0.1:9000;\n"
            conf += "\t\t\t\tfastcgi_param SCRIPT_FILENAME $document_root%s;\n" % scriptName
            conf += "\t\t\t\tinclude fastcgi_params;\n"
            return conf
        for path in locations:
            appNginxConf += "\t\tlocation %s {\n" % path
            # == ROOT
            root = locations[path].get("root", "") or ""
            appNginxConf += "\t\t\troot \"%s\";\n" % (
                ("%s/%s" % (self.APPLICATION_DIRECTORY, root.strip("/"))).rstrip("/")
            )
            # == HEADERS
            headers = locations[path].get("headers", {})
            for headerName in headers:
                appNginxConf += "\t\t\tadd_header %s %s;\n" % (
                    headerName,
                    headers[headerName]
                )
            # == PASSTHRU
            passthru = locations[path].get("passthru", False)
            if passthru and not locations[path].get("scripts", False):
                if passthru == True: passthru = "/index.php"
                appNginxConf += "\t\t\tlocation ~ /%s {\n" % passthru.strip("/")
                appNginxConf += "\t\t\t\tallow all;\n"
                appNginxConf += addFastCgi(passthru)
                appNginxConf += "\t\t\t}\n"
                #appNginxConf += "\t\tlocation / {\n"
                appNginxConf += "\t\t\ttry_files $uri /%s$is_args$args;\n" % passthru.strip("/")
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
            appNginxConf += "\t\t}\n"
        return appNginxConf

    def start(self):
        BasePlatformApplication.start(self)
        container = self.getContainer()
        if not container: return
        # link php.ini in app root
        self.logger.info(
            "Compile php.ini."
        )
        self.runCommand(
            "[ -f /app/php.ini ] && ln -s /app/php.ini /usr/local/etc/php/conf.d/02-app.ini || true"
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
            self.stop()
            return self.start()
        # nginx config
        nginxConfFileObj = io.BytesIO(
            bytes(str(self.generateNginxConfig()).encode("utf-8"))
        )
        self.uploadFile(
            nginxConfFileObj,
            "/etc/nginx/app.conf"
        )
        # start nginx + other services
        self.logger.info(
            "Start Nginx."
        )
        self.runCommand(
            """
            service nginx start
            """
        )
