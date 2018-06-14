import os
import json
import base64
import docker
import logging
import io
from platform_cc.container import Container
from platform_cc.parser.routes import RoutesParser
from platform_cc.exception.state_error import StateError

class BasePlatformApplication(Container):

    """
    Base class for managing Platform.sh applications.
    """

    """
    Directory inside container to mount application to.
    """
    APPLICATION_DIRECTORY = "/app"

    """
    Directory inside container to mount storage to.
    """
    STORAGE_DIRECTORY = "/mnt/storage"

    def __init__(self, project, config):
        """
        Constructor.

        :param project: Project data
        :param config: Service configuration
        """
        self.config = dict(config)
        Container.__init__(
            self,
            project,
            self.config.get(
                "name",
                ""
            )
        )
        self.logger = logging.getLogger(
            "%s.%s.%s" % (
                __name__,
                self.project.get("short_uid"),
                self.getName()
            )
        )

    def getContainerVolumes(self):
        return {
            self.project.get("path") : {
                "bind" : self.APPLICATION_DIRECTORY,
                "mode" : "rw"                
            },
            self.getVolumeName() : {
                "bind" : self.STORAGE_DIRECTORY,
                "mode" : "rw"
            }
        }

    def getContainerEnvironmentVariables(self):
        # get platform relationships
        platformRelationships = {}
        for key, value in self.config.get("relationships", {}).items():
            value = value.strip().split(":")
            platformRelationships[key] = [
                self.project.get("services", {})
                    .get(value[0], {})
                    .get("platform_relationships", {})
                    .get(value[1])
            ]
        routesParser = RoutesParser(self.project)
        # get subnet from project network, used for trusted proxy
        network = self.getNetwork()
        trustedProxies = "%s,127.0.0.1" % (
            str(network.attrs.get("IPAM", {}).get("Config",[{}])[0].get("Subnet"))
        )
        # sest env vars
        envVars = {
            "PLATFORM_APP_DIR"          : self.APPLICATION_DIRECTORY,
            "PLATFORM_APPLICATION"      : "",
            "PLATFORM_APPLICATION_NAME" : self.getName(),
            "PLATFORM_BRANCH"           : "",
            "PLATFORM_DOCUMENT_ROOT"    : "/",
            "PLATFORM_ENVIRONMENT"      : "",
            "PLATFORM_PROJECT"          : self.project.get("uid", ""),
            "PLATFORM_RELATIONSHIPS"    : base64.b64encode(
                bytes(str(json.dumps(platformRelationships)).encode("utf-8"))
            ).decode("utf-8"),
            "PLATFORM_ROUTES"           : base64.b64encode(
                bytes(str(json.dumps(routesParser.getRoutesEnvironmentVariable())).encode("utf-8"))
            ),
            "PLATFORM_TREE_ID"          : "",
            "PLATFORM_VARIABLES"        : base64.b64encode(
                bytes(str(json.dumps(self.project.get("variables", {}))).encode("utf-8"))
            ).decode("utf-8"),
            "PLATFORM_PROJECT_ENTROPY"  : self.project.get("entropy", ""),
            "TRUSTED_PROXIES"           : trustedProxies
        }
        # set env vars from app variables
        for key, value in self.config.get("variables", {}).get("env", {}).items():
            envVars[key.strip().upper()] = str(value)
        # set env vars from project variables
        for key, value in self.project.get("variables", {}).items():
            if not key.startswith("env:"): continue
            key = key[4:]
            envVars[key.strip().upper()] = str(value)
        
        return envVars

    def getContainerWorkingDirectory(self):
        return self.APPLICATION_DIRECTORY

    def getType(self):
        """
        Get application type.

        :return: Application type
        :rtype: str
        """        
        return self.config.get("type")

    def setupMounts(self):
        """
        Setup application defined mount points.
        """
        # project option 'use_mount_volumes' must be true
        configMounts = self.config.get("mounts", {})
        self.logger.info(
            "Found %s mount point(s).",
            len(configMounts)
        )
        for mountDest, config in configMounts.items():
            mountSrc = ""
            if type(config) is dict:
                if not config.get("source") == "local": continue
                mountSrc = config.get("source_path", "").strip("/")
            elif type(config) is str:
                localMountPrefx = "shared:files/"
                if not config.startswith(localMountPrefx): continue
                mountSrc = config[len(localMountPrefx):].strip("/")
            else:
                continue
            self.logger.debug(
                "Bind mount point '%s.'.",
                mountSrc
            )
            mountSrc = os.path.join(
                self.STORAGE_DIRECTORY,
                mountSrc.strip("/")
            )
            mountDest = os.path.join(
                self.APPLICATION_DIRECTORY,
                mountDest.strip("/")
            )
            self.runCommand(
                "mkdir -p %s && mkdir -p %s" % (
                    mountSrc,
                    mountDest,
                ),
                "root"
            )
            if not self.project.get("config", {}).get("option_use_mount_volumes"): continue
            self.runCommand(
                "mount -o user_xattr --bind %s %s" % (
                    mountSrc,
                    mountDest
                ),
                "root"
            )

    def installSsh(self):
        """
        Install SSH key and known hosts file.
        """
        sshDatas = [
            ["ssh_key", "/app/.ssh/id_rsa"],
            ["ssh_known_hosts", "/app/.ssh/known_hosts"]
        ]
        self.runCommand(
            "mkdir -p /app/.ssh && chown -f -R web:web /app/.ssh"
        )
        for sshData in sshDatas:
            data = self.project.get("config", {}).get(sshData[0])
            if not data: continue
            self.logger.info(
                "Install '%s' from project config." % sshData[0]
            )
            data = str(base64.b64decode(data))
            dataFileObject = io.BytesIO(data)
            self.uploadFile(
                dataFileObject,
                "/tmp/.ssh_file" # can't upload file to a mount directory, so upload to tmp and copy
            )
            self.runCommand(
                "mv /tmp/.ssh_file %s && chmod -f 0600 %s" % (
                    sshData[1],
                    sshData[1]
                )
            )

    def installCron(self):
        """
        Install cron tasks and enable cron in application container.
        """
        # cron must be enabled via options
        if not self.project.get("config", {}).get("option_enable_cron"): return
        # create cron directory if not exist
        self.runCommand(
            "mkdir -p /etc/cron.d"
        )
        # itterate crons make cron files
        crons = self.config.get("crons", {})
        self.logger.info(
            "Installing %s cron task(s)." % str(len(crons))
        )
        for name, cron in crons.items():
            spec = cron.get("spec", "*/5 * * * *") # default is every 5 minutes
            cmd = cron.get("cmd", "")
            if not cmd: continue
            self.logger.debug(
                "Installing '%s' cron." % name
            )
            fileObj = io.BytesIO(
                bytes(
                    "%s web %s" % (
                        spec,
                        cmd
                    )
                )
            )
            self.uploadFile(
                fileObj,
                "/etc/cron.d/%s" % name
            )
        # start cron
        self.logger.info("Start cron daemon.")
        self.runCommand("cron")
            
    def build(self):
        """
        Run commands needed to get container ready for given
        application. Also runs build hooks commands.
        """
        self.logger.info(
            "Building application."
        )
        self.installSsh()
        output = self.runCommand(
            self.config.get("hooks", {}).get("build", "")
        )
        # commit container
        self.logger.info(
            "Commit container."
        )
        self.commit()
        return output

    def deploy(self):
        """
        Run deploy hook commands.
        """
        self.logger.info(
            "Run deploy hooks."
        )
        return self.runCommand(
            self.config.get("hooks", {}).get("deploy", ""),
            "web"
        )

    def start(self):
        # ensure all required services are available
        projectServices = self.project.get("services", {})
        serviceNames = list(self.config.get("relationships", {}).values())
        for serviceName in serviceNames:
            serviceName = serviceName.strip().split(":")[0]
            projectService = projectServices.get(serviceName)
            if not projectService or not projectService.get("running"):
                raise StateError(
                    "Application '%s' depends on service '%s' which is not running." % (
                        self.getName(),
                        serviceName
                    )
                )
        Container.start(self)