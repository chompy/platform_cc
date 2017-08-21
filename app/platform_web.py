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

    def generateNginxConfig(self):
        """ Generate nginx config file for application. """
        webConfig = self.app.config.getWeb()
        locations = webConfig.get("locations", {})
        baseDocker = PlatformDocker(
            self.app.config,
            self.app.config.getDockerImage()
        )
        webDocker = PlatformDocker(
            self.app.config,
            self.WEB_DOCKER_IMAGE
        )
        webProvisionConfig = webDocker.getProvisioner().config
        baseNginxConfig = webProvisionConfig.get("nginx.conf", "")
        appNginxConf = ""
        for path in locations:
            appNginxConf += "location %s {\n" % path
            # root
            appNginxConf += "\t\troot \"%s\";\n" % (
                "/app/%s" % (locations[path].get("root", "").strip("/"))
            )
            # optional directives
            for directive in locations[path]:
                # passthru
                if directive == "passthru":
                    value = locations[path][directive].strip()
                    if not value: value = "/index.php"
                    appNginxConf += "\t\tindex %s;\n" % (value.lstrip("/"))
                    appNginxConf += "\t\tlocation ~ [^/]\.php(/|$) {\n"
                    appNginxConf += "\t\t\tfastcgi_split_path_info ^(.+?\.php)(/.*)$;\n"
                    appNginxConf += "\t\t\tfastcgi_pass %s:9000;\n" % (baseDocker.containerId)
                    appNginxConf += "\t\t\tfastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;\n"
                    appNginxConf += "\t\t\tfastcgi_index %s;\n" % (value.lstrip("/"))
                    appNginxConf += "\t\t\tinclude fastcgi_params;\n"
                    appNginxConf += "\t\t}\n"
            appNginxConf += "\t}\n"
        # TODO rules, expire, scripts, allow, headers
        return baseNginxConfig.replace("{{APP_WEB}}", appNginxConf)

    def start(self):
        """ Start app web handler. """
        print_stdout("> Starting web handler for '%s' app." % self.app.config.getName())
        docker = PlatformDocker(
            self.app.config,
            self.WEB_DOCKER_IMAGE
        )
        docker.start()
        docker.getProvisioner().copyStringToFile(
            self.generateNginxConfig(),
            "/etc/nginx/nginx.conf"
        )
        docker.getContainer().restart()

    def stop(self):
        """ Stop app web handler. """
        print_stdout("> Stopping web handler for '%s' app." % self.app.config.getName())
        docker = PlatformDocker(
            self.app.config,
            self.WEB_DOCKER_IMAGE
        )
        docker.stop()