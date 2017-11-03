import os
import yaml
from config.platform_config import PlatformConfig
from platform_docker import PlatformDocker

class PlatformWeb:

    """ Provide web access to app via nginx docker container. """

    WEB_DOCKER_IMAGE = "nginx:1.13"

    def __init__(self, app):
        self.app = app
        self.logger = self.app.logger
        self.docker = PlatformDocker(
            self.app.config,
            "%s_web" % self.app.config.getName(),
            self.WEB_DOCKER_IMAGE,
            self.logger
        )
        self.logIndent = 1
        self.docker.logIndent = self.logIndent + 1

    def generateNginxConfig(self):
        """ Generate nginx config file for application. """
        baseNginxConfig = self.docker.getProvisioner().config.get("web_conf", "")
        return baseNginxConfig.replace(\
            "{{APP_WEB}}",
            self.app.docker.getProvisioner().generateNginxConfig()
        )

    def start(self):
        """ Start app web handler. """
        if self.logger:
            self.logger.logEvent(
                "Starting web server.",
                self.logIndent
            )
        self.docker.start()
        self.docker.getProvisioner().copyStringToFile(
            str(self.generateNginxConfig()),
            "/etc/nginx/nginx.conf"
        )
        self.docker.getContainer().restart()

    def stop(self):
        """ Stop app web handler. """
        if self.logger:
            self.logger.logEvent(
                "Stopping web server.",
                self.logIndent
            )
        self.docker.stop()
