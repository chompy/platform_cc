from __future__ import absolute_import
import os
import yaml
import docker
from .config.platform_config import PlatformConfig
from .platform_docker import PlatformDocker

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
        return baseNginxConfig.replace(
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

        # provision if needed
        commitImage = "%s:%s" % (self.docker.DOCKER_COMMIT_REPO, self.docker.getTag())
        try:
            image = self.docker.dockerClient.images.get(commitImage)
        except docker.errors.NotFound:
            image = None
        if not image:
            self.docker.provision()

        # update conf files
        self.docker.getProvisioner().copyStringToFile(
            self.docker.getProvisioner().config.get("web_conf", ""),
            "/etc/nginx/nginx.conf"
        )
        self.docker.getProvisioner().copyStringToFile(
            self.app.docker.getProvisioner().generateNginxConfig(),
            "/etc/nginx/app.conf"
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