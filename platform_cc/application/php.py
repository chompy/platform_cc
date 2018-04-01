import io
import os
import json
from .base import BasePlatformApplication
from exception.container_command_error import ContainerCommandError

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
        # create 'web' user
        self.runCommand(
            "useradd -d /app -m -p secret~ --uid %s web || true" % (
                self.project.get("config", {}).get("web_uid", self.DEFAULT_WEB_UID)
            )            
        )
        # provision container
        self.runCommand(
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
            echo "date.timezone = UTC" > /usr/local/etc/php/conf.d/main.ini
            echo "memory_limit = 512M" >> /usr/local/etc/php/conf.d/main.ini
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
        # install extensions
        extInstall = self.config.get("runtime", {}).get("extensions", [])
        for extension in extInstall:
            if type(extension) is not str: continue
            command = self.getExtensionInstallCommand(extension)
            if not command: continue
            self.runCommand(command)
        # build hooks
        try:
            self.runCommand(
                self.config.get("hooks", {}).get("build", ""),
                "web"
            )
        # allow build hooks to fail...for now
        except ContainerCommandError:
            pass
        # clean up
        self.runCommand(
            """
            apt-get clean
            """
        )
        # commit container
        container = self.getContainer()
        container.commit(
            self.COMMIT_REPOSITORY_NAME,
            "%s_%s" % (
                self.getName(),
                self.project.get("short_uid")
            )
        )

    def start(self):
        BasePlatformApplication.start(self)
        container = self.getContainer()
        if not container: return
        # setup mount points
        self.setupMounts()
        # link php.ini in app root
        self.runCommand(
            "[ -f /app/php.ini ] && ln -s /app/php.ini /usr/local/etc/php/conf.d/app.ini || true"
        )
        # build php.ini from config vars
        phpIniConfig = self.config.get("variables", {}).get("php", {})
        phpIniFileObj = io.StringIO()
        for key, value in phpIniConfig.items():
            phpIniFileObj.write("%s = %s\n" % (key, value))
        for key, value in self.project.get("variables", {}).items():
            if not key or key[0:4] != "php:": continue
            phpIniFileObj.write("%s = %s\n" % (key[4:], value))
        self.uploadFile(
            phpIniFileObj,
            "/usr/local/etc/php/conf.d/app2.ini"
        )
        # restart container to reload conf changes
        container.restart()
        # not yet built/provisioned
        if self.getDockerImage() == self.getBaseImage():
            self.build()
            container.restart()
        # start nginx + other services
        self.runCommand(
            """
            service nginx start
            """
        )
