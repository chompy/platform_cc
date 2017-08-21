import os
import difflib
import io
import hashlib
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a PHP container. """

    PHP_EXTENSION_DEPENDENCIES = {
        "memcached" : {
            "php:5.6" : """
                apt-get install -y libmemcached-dev zlib1g-dev
                print '\n' | pecl install memcached-2.2.0
                docker-php-ext-enable memcached
            """,
            "php:7" : """
                apt-get install -y libmemcached-dev zlib1g-dev
                print '\n' | pecl install memcached-2.2.0
                docker-php-ext-enable memcached
            """,
        },
        "gd" : {
            "php:5.6" : """
                apt-get install -y libfreetype6-dev libjpeg62-turbo-dev libpng12-dev
                docker-php-ext-configure --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/
                docker-php-ext-install -j$(nproc) gd
            """
        },
        "imagick" : {
            "php:5.6" : """
                apt-get install -y imagemagick libmagickcore-dev libmagickwand-dev
                print '\n' | pecl install imagick
                docker-php-ext-enable imagick
            """
        },
        "intl" : {
            "php:5.6" : """
                apt-get install -y libicu-dev
                docker-php-ext-install -j$(nproc) intl            
            """
        },
        "pdo_mysql" : {
            "php:5.6" : """
                docker-php-ext-install -j$(nproc) pdo_mysql
            """
        },
        "mcrypt" : {
            "php:5.6" : """
                apt-get install -y libmcrypt-dev
                docker-php-ext-install -j$(nproc) mcyrpt
            """
        },
        "xsl" : {
            "php:5.6" : """
                apt-get install -y libxslt1-dev
                docker-php-ext-install -j$(nproc) xsl
            """
        },
        "zip" : {
            "php:5.6" : """
                docker-php-ext-install -j$(nproc) zip
            """
        }
    }

    PHP_EXTENSION_CORE = ["curl", "json", "sqlite3"]

    def provision(self):
        # add 'web' user
        print_stdout("  - Create 'web' user...", False)
        password = self.randomString(10)
        self.container.exec_run(
            ["useradd", "-d", "/app", "-m", "-p", password, "web"]
        )
        print_stdout("done.")

        # install additional dependencies
        print_stdout("  - Install additional dependencies...", False)
        self.container.exec_run(
            ["apt-get", "update"]
        )
        self.container.exec_run(
            ["apt-get", "install", "-y", "rsync", "git", "unzip"]
        )
        print_stdout("done.")

        # php conf file configure
        print_stdout("  - Additional PHP configuration...", False)
        self.container.exec_run(
            ["sed", "-i", "s/user = .*/user = web/g", "/usr/local/etc/php-fpm.d/www.conf"]
        )
        self.container.exec_run(
            ["sed", "-i", "s/group = .*/group = web/g", "/usr/local/etc/php-fpm.d/www.conf"]
        )
        self.copyStringToFile(
            "date.timezone = UTC",
            "/usr/local/etc/php/conf.d/timezone.ini"
        )
        print_stdout("done.")

        # composer install
        if self.platformConfig.getBuildFlavor() == "composer":
            print_stdout("  - Install composer...", False)
            self.container.exec_run(
                ["php", "-r", "copy('https://getcomposer.org/installer', 'composer-setup.php');"]
            )
            self.container.exec_run(
                ["php", "composer-setup.php", "--install-dir=/usr/local/bin"]
            )
            self.container.exec_run(
                ["rm", "composer-setup.php"]
            )
            print_stdout("done.")

        # install extensions
        print_stdout("  - Install extensions...")
        extensions = self.platformConfig.getRuntime().get("extensions", [])
        for extension in extensions:
            print_stdout("    - %s..." % (extension), False)
            if extension in self.PHP_EXTENSION_CORE:
                print_stdout("already installed (core extension).")
                continue
            if extension not in self.PHP_EXTENSION_DEPENDENCIES:
                print_stdout("not available.")
                continue
            extensionDeps = self.PHP_EXTENSION_DEPENDENCIES[extension]
            depCmdKey = difflib.get_close_matches(
                self.image,
                extensionDeps.keys(),
                1
            )
            if not depCmdKey:
                print_stdout("not available.")
                continue
            self.container.exec_run(
                ["sh", "-c", extensionDeps[depCmdKey[0]]]
            )
            print_stdout("done.")

    def preBuild(self):
        # rsync app
        print_stdout("  - Copy application to container...", False)
        self.container.exec_run(
            ["rsync", "-a", "--exclude", ".platform", "--exclude", ".git", "--exclude", ".platform.app.yaml", "/mnt/app/", "/app"]
        )
        self.container.exec_run(
            ["chown", "-R", "web:web", "/app"]
        )
        print_stdout("done.")

        # install ssh key
        print_stdout("  - Install SSH key file...", False)
        self.container.exec_run(
            ["mkdir", "-p", "/app/.ssh"]
        )
        self.copyFile(
            os.path.join(
                self.platformConfig.getDataPath(),
                "id_rsa"
            ),
            "/app/.ssh/id_rsa"
        )
        self.container.exec_run(
            ["chmod", "0600", "/app/.ssh/id_rsa"]
        )
        if os.path.exists(os.path.join(self.platformConfig.getDataPath(), "known_hosts")):
            self.copyFile(
                os.path.join(
                    self.platformConfig.getDataPath(),
                    "known_hosts"
                ),
                "/app/.ssh/known_hosts"
            ) 
        print_stdout("done.")

        # run 'composer install'
        if self.platformConfig.getBuildFlavor() == "composer":
            print_stdout("  - Running composer...", False)
            self.container.exec_run(
                ["php", "-d", "memory_limit=-1", "/usr/local/bin/composer.phar", "install", "-n", "-d", "/app"],
                user="web"
            )
            print_stdout("done.")

    def getUid(self):
        """ Generate unique id based on configuration. """
        hashStr = self.image
        hashStr += str(self.platformConfig.getBuildFlavor())
        extensions = self.platformConfig.getRuntime().get("extensions", [])
        for extension in extensions:
            if extension not in self.PHP_EXTENSION_CORE and extension in self.PHP_EXTENSION_DEPENDENCIES:
                hashStr += extension
        return hashlib.sha256(hashStr).hexdigest()