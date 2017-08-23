from platform_service_config import PlatformServiceConfig
from platform_docker import PlatformDocker
from app.platform_utils import log_stdout

class PlatformService:

    """ Platform service handler. """

    def __init__(self, appConfig, name):
        self.name = str(name).strip()
        self.appConfig = appConfig
        self.config = PlatformServiceConfig(
            self.appConfig.projectHash,
            self.appConfig.appPath,
            name
        )
        self.docker = PlatformDocker(
            self.config,
            "%s_%s" % (appConfig.getName(), self.config.getName())
        )

    def start(self):
        """ Start service. """
        log_stdout("Starting '%s' service." % self.config.getName())
        if not self.config.getDockerImage():
            log_stdout("No docker image available, skipping", 1)
            return
        self.docker.start()

    def stop(self):
        """ Stop service. """
        log_stdout("Stopping '%s' service." % self.config.getName())
        if not self.config.getDockerImage():
            log_stdout("No docker image available, skipping", 1)
            return
        self.docker.stop()

    def getRelationShipType(self):
        return self.docker.getProvisioner().getRelationshipType()

    def getRelationships(self):
        """ Get service relationships to expose to application. """
        return self.docker.getProvisioner().getRelationships()