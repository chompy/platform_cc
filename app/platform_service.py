from platform_service_config import PlatformServiceConfig
from platform_docker import PlatformDocker
from app.platform_utils import log_stdout

class PlatformService:

    """ Platform service definition. """

    PLATFORM_SERVICE_DOCKER_IMAGES = {
        "mysql":                   "mariadb:10.2",
        "mysql:10.2":              "mariadb:10.2",
        "mysql:10.1":              "mariadb:10.1",
        "mysql:10.0":              "mariadb:10.0",
        "mysql:5.5":               "mariadb:5.5",
        "memcached":               "memcached:1",
        "memcached:1.4":           "memcached:1"
    }

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
        log_stdout("Stopping '%s' service." % self.config.getName())
        if not self.config.getDockerImage():
            log_stdout("No docker image available, skipping", 1)
            return
        self.docker.stop()