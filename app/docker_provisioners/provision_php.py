import os
import difflib
import io
from provision_base import DockerProvisionBase

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
                apt-get install -y libmagickcore-dev libmagickwand-dev
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
        print "  - Create 'web' user...",
        password = self.randomString(10)
        self.container.exec_run(
            ["useradd", "-d", "/app", "-m", "-p", password, "web"]
        )
        print "done."

        # install additional dependencies
        print "  - Install additional dependencies...",
        self.container.exec_run(
            ["apt-get", "update"]
        )
        self.container.exec_run(
            ["apt-get", "install", "-y", "rsync", "git", "unzip"]
        )
        print "done."

        # rsync app
        print "  - Copy application to container...",
        self.container.exec_run(
            ["rsync", "-a", "--exclude", ".platform", "--exclude", ".git", "--exclude", ".platform.app.yaml", "/mnt/app/", "/app"]
        )
        print "done."

        # install ssh key
        print "  - Install SSH key file...",
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
        print "done."

        # composer install
        print "  - Install composer...",
        self.container.exec_run(
            ["php", "-r", "copy('https://getcomposer.org/installer', 'composer-setup.php');"]
        )
        self.container.exec_run(
            ["php", "composer-setup.php", "--install-dir=/usr/local/bin"]
        )
        self.container.exec_run(
            ["rm", "composer-setup.php"]
        )
        print "done."

        # install extensions
        print "  - Install extensions..."
        extensions = self.platformConfig.getRuntime().get("extensions", [])
        for extension in extensions:
            print "    - %s..." % (extension),
            if extension in self.PHP_EXTENSION_CORE:
                print "already installed (core extension)."
                continue
            if extension not in self.PHP_EXTENSION_DEPENDENCIES:
                print "not available."
                continue
            extensionDeps = self.PHP_EXTENSION_DEPENDENCIES[extension]
            depCmdKey = difflib.get_close_matches(
                self.image,
                extensionDeps.keys(),
                1  
            )
            if not depCmdKey:
                print "not available."
                continue
            self.copyStringToFile(
                extensionDeps[depCmdKey[0]],
                "/php_install_ext"
            )
            self.container.exec_run(
                ["chmod", "+x", "/php_install_ext"]
            )
            self.container.exec_run(
                ["sh", "/php_install_ext"]
            )
            self.container.exec_run(
                ["rm", "/php_install_ext"]
            )
            print "done."