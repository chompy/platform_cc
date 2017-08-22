import os
import yaml
from platform_config import PlatformConfig
from platform_docker import PlatformDocker
from platform_utils import print_stdout

class PlatformWeb:

    """ Provide web access to app via nginx docker container. """

    WEB_DOCKER_IMAGE = "nginx:1.13"

    def __init__(self, app):
        self.app = app
        self.docker = PlatformDocker(
            self.app.config,
            self.WEB_DOCKER_IMAGE,
            "web"
        )

    def getDocker(self):
        return Pl

    def generateNginxConfig(self):
        """ Generate nginx config file for application. """
        webConfig = self.app.config.getWeb()
        locations = webConfig.get("locations", {})

        webProvisionConfig = self.docker.getProvisioner().config
        baseNginxConfig = webProvisionConfig.get("nginx.conf", "")
        appNginxConf = ""
        for path in locations:

            appNginxConf += "location %s {\n" % path
            
            # root
            appNginxConf += "\t\troot \"%s\";\n" % (
                "/app/%s" % (locations[path].get("root", "").strip("/"))
            )

            # passthru
            passthru = False
            if "passthru" in locations[path]:
                passthru = locations[path]["passthru"].strip()
                if not passthru: passthru = "/index.php"
                appNginxConf += "\t\tlocation ~ /%s {\n" % (passthru.strip("/"))
                appNginxConf += "\t\t\tfastcgi_split_path_info ^(.+?\.php)(/.*)$;\n"
                appNginxConf += "\t\t\tfastcgi_pass %s:9000;\n" % (self.app.docker.containerId)
                appNginxConf += "\t\t\tfastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;\n"
                appNginxConf += "\t\t\tfastcgi_index %s;\n" % (passthru.lstrip("/"))
                appNginxConf += "\t\t\tinclude fastcgi_params;\n"
                appNginxConf += "\t\t}\n"
                appNginxConf += "\t\tlocation / {\n"
                if passthru:
                    appNginxConf += "\t\t\ttry_files $uri /%s$is_args$args;\n" % (passthru.strip("/"))
                appNginxConf += "\t\t}\n"

            # scripts
            scripts = False
            if "scripts" in locations[path]:
                scripts = locations[path]["scripts"]
            appNginxConf += "\t\tlocation ~ [^/]\.php(/|$) {\n"
            if scripts:
                appNginxConf += "\t\t\tfastcgi_split_path_info ^(.+?\.php)(/.*)$;\n"
                appNginxConf += "\t\t\tfastcgi_pass %s:9000;\n" % (self.app.docker.containerId)
                appNginxConf += "\t\t\tfastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;\n"
                if passthru:
                    appNginxConf += "\t\t\tfastcgi_index %s;\n" % (passthru.lstrip("/"))
                appNginxConf += "\t\t\tinclude fastcgi_params;\n"
            else:
                appNginxConf += "\t\t\tdeny all;\n"
            appNginxConf += "\t\t}\n"
            appNginxConf += "\t}\n"

        # TODO rules, expire, allow, headers
        return baseNginxConfig.replace("{{APP_WEB}}", appNginxConf)

    def start(self):
        """ Start app web handler. """
        print_stdout("> Starting web handler for '%s' app." % self.app.config.getName())
        self.docker.start()
        self.docker.getProvisioner().copyStringToFile(
            self.generateNginxConfig(),
            "/etc/nginx/nginx.conf"
        )
        self.docker.getContainer().restart()

    def stop(self):
        """ Stop app web handler. """
        print_stdout("> Stopping web handler for '%s' app." % self.app.config.getName())
        self.docker.stop()