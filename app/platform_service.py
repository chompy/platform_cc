from config.platform_service_config import PlatformServiceConfig
from platform_docker import PlatformDocker

class PlatformService:

    """ Platform service handler. """

    def __init__(self, appConfig, name, logger = None):
        self.name = str(name).strip()
        self.appConfig = appConfig
        self.config = PlatformServiceConfig(
            self.appConfig.projectHash,
            self.appConfig.appPath,
            name
        )
        self.docker = PlatformDocker(
            self.config,
            "%s_%s" % (appConfig.getName(), self.config.getName()),
            self.config.getDockerImage(),
            logger
        )
        self.docker.logIndent += 1
        self.logger = logger
        self.logIndent = 1

    def start(self):
        """ Start service. """
        if self.logger:
            self.logger.logEvent(
                "Starting '%s' service..." % self.config.getName(),
                self.logIndent
            )
        if not self.config.getDockerImage():
            if self.logger:
                self.logger.logEvent(
                    "No docker image available, skipped.",
                    self.logIndent + 1
                )
            return
        self.docker.start()

    def stop(self):
        """ Stop service. """
        if self.logger:
            self.logger.logEvent(
                "Stopping '%s' service." % self.config.getName(),
                self.logIndent
            )
        if not self.config.getDockerImage():
            if self.logger:
                self.logger.logEvent(
                    "No docker image available, skipped.",
                    self.logIndent + 1
                )
            return
        self.docker.stop()
