"""
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
"""

import os
import urllib
import json
import ast
import base64
import logging
from platform_cc.container import Container
from platform_cc.application.php import PhpApplication
from .config import PlatformShConfig
from .api import PlatformShApi
from .fetch import getPlatformShFetcher
from platform_cc.project import PlatformProject
from .exception.api_error import PlatformShApiError
from platform_cc.exception.container_command_error import ContainerCommandError
from platform_cc.exception.state_error import StateError

class PlatformShCloner(Container):

    """ Tools to clone a Platform.sh environment. """

    def __init__(self, projectId, environment, path, config = None, dockerClient = None):
        if not config:
            config = PlatformShConfig()
        self.config = config
        Container.__init__(
            self,
            project={
                "_path" : str(path),
                "_project_id" : str(projectId),
                "_environment" : str(environment),
                "uid" : "__psh_clone",
                "short_uid" : "__psh_clone",
                "config" : {},
            },
            name="psh_clone",
            dockerClient=dockerClient
        )
        self.logger = logging.getLogger(
            __name__
        )
        self._pshProject = None
        self._pshDeploy = None
        self._pshEnv = None

    def _getPshProject(self):
        if self._pshProject:
            return self._pshProject
        pshApi = PlatformShApi(self.config)
        self._pshProject = pshApi.getProjectInfo(self.project.get("_project_id"))
        return self._pshProject

    def _getPshEnv(self):
        if self._pshEnv:
            return self._pshEnv
        pshApi = PlatformShApi(self.config)
        self._pshEnv = pshApi.getEnvironmentInfo(self.project.get("_project_id"), self.project.get("_environment"))
        return self._pshEnv
    
    def _getPshDeploy(self):
        if self._pshDeploy:
            return self._pshDeploy
        pshApi = PlatformShApi(self.config)
        self._pshDeploy = pshApi.getDeployment(self.project.get("_project_id"), self.project.get("_environment"))
        return self._pshDeploy

    def getBaseImage(self):
        # just use php as the image contains ssh
        return PhpApplication.DOCKER_IMAGE_MAP["php"]

    def getDockerImage(self):
        return self.getBaseImage()

    def getCommitImage(self):
        return None

    def getContainerCommand(self):
        return "sh"

    def getContainerEnvironmentVariables(self):
        pshApi = PlatformShApi(self.config)
        return {
            "PSH_ACCESS_TOKEN" : self.config.getAccessToken(),
            "PSH_SSH_KEY" : pshApi.getSshKeypair()[1]
        }

    def getContainerVolumes(self):
        return {
            os.path.abspath(self.project.get("_path")) : {
                "bind": PhpApplication.APPLICATION_DIRECTORY,
                "mode": "rw"
            }
        }

    def start(self):

        # get git url to use to add ssh known host
        pshProject = self._getPshProject()
        pshGitUrl = pshProject.get("repository", {}).get("url", "")
        if not pshGitUrl:
            raise PlatformShApiError("Could not fetch Git URL for project '%s'" % self.project.get("_project_id"))
        parsedGitUrl = urllib.parse.urlparse("ssh://%s" % pshGitUrl)

        # retrieve ssh url
        pshEnv = self._getPshEnv()
        sshUrl = pshEnv.get("_links", {}).get("ssh", {}).get("href")
        parsedSshUrl = urllib.parse.urlparse(sshUrl)
        sshUrl = sshUrl.replace("ssh://", "")

        # start container
        Container.start(self)

        # get current user id
        try:
            currentUserId = os.getuid()
        except AttributeError:
            currentUserId = 1000
        if currentUserId <= 0:
            currentUserId = 1000

        # setup ssh + update user id
        # TODO all custom entry in known_hosts
        cmd = """
        mkdir -p ~/.ssh
        echo "$PSH_SSH_KEY" > ~/.ssh/id_rsa
        chmod 0600 ~/.ssh/id_rsa
        ssh-keyscan %s > ~/.ssh/known_hosts
        ssh-keyscan %s >> ~/.ssh/known_hosts
        ssh-keyscan gitlab.com >> ~/.ssh/known_hosts
        ssh-keyscan github.com >> ~/.ssh/known_hosts
        chmod 0600 ~/.ssh/known_hosts
        usermod -u %s web
        """ % (
            parsedSshUrl.hostname,
            parsedGitUrl.hostname,
            currentUserId
        )
        try:
            self.runCommand(cmd)
        except Exception as e:
            self.logger.error("An error occured, stopping container...")
            self.stop()
            raise e

    def clone(self, skipMountSync=False, skipServiceSync=False):
        """ 
        Clone platform.sh environment by using its REST API, git clone
        and accessing environment via ssh for asset dumps.
        """
        
        self.logger.info("Retrieve project information (%s)." % self.project.get("_project_id"))
        pshApi = PlatformShApi(self.config)
        
        # get project info to retrieve git URL
        pshProject = pshApi.getProjectInfo(self.project.get("_project_id"))
        pshGitUrl = pshProject.get("repository", {}).get("url", "")
        if not pshGitUrl:
            raise PlatformShApiError("Could not fetch Git URL for project '%s'" % self.project.get("_project_id"))
        parsedGitUrl = urllib.parse.urlparse("ssh://%s" % pshGitUrl)
        
        # retrieve ssh url
        pshEnv = self._getPshEnv()
        sshUrl = pshEnv.get("_links", {}).get("ssh", {}).get("href")
        parsedSshUrl = urllib.parse.urlparse(sshUrl)
        sshUrl = sshUrl.replace("ssh://", "")

        # start clone container
        self.start()

        # git clone
        if not os.path.exists(os.path.join(self.project.get("_path"), self.project.get("_project_id"))):
            self.logger.info(
                "Cloning project %s (%s)." % (
                    self.project.get("_project_id"),
                    self.project.get("_environment")
                )
            )
            try:
                self.runCommand(
                    """
                    mkdir -p %s
                    cd %s
                    git clone "%s" --recursive --branch "%s"
                    chown -R web:web %s
                    chmod -R g+rw %s
                    """ % (
                        PhpApplication.APPLICATION_DIRECTORY,
                        PhpApplication.APPLICATION_DIRECTORY,
                        pshGitUrl,
                        self.project.get("_environment"),
                        self.project.get("_project_id"),
                        self.project.get("_project_id")
                    )                    
                )
            except Exception as e:
                self.logger.error("An error occured, stopping container...")
                self.stop()
                raise e

        # get PCC project
        pccProject = PlatformProject.fromPath(
            os.path.join(
                self.project.get("_path"),
                self.project.get("_project_id")
            )
        )

        # add ssh key to project
        pccProject.config.set(
            "ssh_key",  
            base64.b64encode(
                bytes(
                    str(
                        pshApi.getSshKeypair()[1]
                    ).encode("utf-8")
                )
            ).decode("utf-8")
        )

        # set vars
        self.syncVars(pccProject)

        # CUSTOM contextual code vars 
        # TODO instead of hard coding maybe add a config file to override vars?
        pccProject.variables.set("env:BUSINESS_HOURS_IGNORE", "1")

        # CUSTOM contextual code behavior, use 'env:PRIMARY_REPO' to change the project
        # dirname and update the git repo remote
        primaryRepo = pccProject.variables.get("env:PRIMARY_REPO")
        if primaryRepo:
            # rename to project name based on primary repo
            projDirName = os.path.splitext(os.path.basename(primaryRepo))[0]
            projFullPath = os.path.join(
                self.project.get("_path"),
                projDirName
            )
            i = 2
            while os.path.exists(projFullPath):
                projDirName = "%s-%d" % (projDirName, i)
                projFullPath = os.path.join(
                    self.project.get("_path"),
                    projDirName
                )
                i += 1
            self.logger.info("Rename project directory to '%s.'" % projDirName)
            self.runCommand(
                """
                mv %s %s
                """ % (
                    os.path.join(PhpApplication.APPLICATION_DIRECTORY, self.project.get("_project_id")),
                    os.path.join(PhpApplication.APPLICATION_DIRECTORY, projDirName),
                )
            )
            # update git remote
            self.logger.info("Updated Git remote to '%s.'" % primaryRepo)
            try:
                self.runCommand(
                    """
                    mv %s %s
                    cd %s
                    git remote set-url origin %s
                    git remote add platform %s || true
                    git remote set-url platform %s
                    """ % (
                        os.path.join(PhpApplication.APPLICATION_DIRECTORY, self.project.get("_project_id")),
                        os.path.join(PhpApplication.APPLICATION_DIRECTORY, projDirName),
                        os.path.join(PhpApplication.APPLICATION_DIRECTORY, projDirName),
                        primaryRepo,
                        sshUrl,
                        sshUrl
                    )
                )
            except Exception as e:
                pass
            pccProject = PlatformProject.fromPath(projFullPath)

        # set psh project id to env var
        pccProject.variables.set("env:PSH_PROJECT_ID", self.project.get("_project_id"))

        # start project
        self.logger.info("Start project.")
        try:
            pccProject.start()
        except Exception as e:
            self.logger.error("An error occured, stopping container...")
            self.stop()
            pccProject.stop()
            raise e

        # add known hosts
        try:
            pccProject.getApplication().runCommand(
                """
                mkdir -p ~/.ssh
                ssh-keyscan %s > ~/.ssh/known_hosts
                ssh-keyscan %s >> ~/.ssh/known_hosts
                ssh-keyscan gitlab.com >> ~/.ssh/known_hosts
                ssh-keyscan github.com >> ~/.ssh/known_hosts            
                chmod 0600 ~/.ssh/known_hosts
                """ % (
                    parsedSshUrl.hostname,
                    parsedGitUrl.hostname,
                ),
                "web"
            )
        except Exception as e:
            self.logger.error("An error occured, stopping container...")
            self.stop()
            pccProject.stop()
            raise e

        # rsync mounts
        if not skipMountSync:
            self.syncMounts(pccProject)
        # sync services
        if not skipServiceSync:
            try:
                self.syncServices(pccProject)
            except Exception as e:
                self.logger.error("An error occured, stopping container...")
                self.stop()
                pccProject.stop()
                raise e

        self.stop()
        pccProject.stop()

    def syncVars(self, pccProject):
        """ Sync project vars with Platform.sh. """
        pshDeploy = self._getPshDeploy()
        for var in pshDeploy[0].get("variables", []):
            key = var.get("name")
            value = var.get("value", "")
            sensitive = var.get("is_sensitive", False)
            if not key: continue
            if sensitive:
                self.logger.info("SKIPPING sensitive project variable '%s.'" % var.get("name"))
                continue
            self.logger.info("Set project variable '%s.'" % var.get("name"))
            try:
                value = ast.literal_eval(value)
                value = json.dumps(value)
            except ValueError:
                pass
            except SyntaxError:
                pass
            pccProject.variables.set(var.get("name"), str(value))

    def syncMounts(self, pccProject):
        """ Sync project mounts with Platform.sh. (Rsync) """
        # get project env info
        pshEnv = self._getPshEnv()
        # retrieve ssh url
        sshUrl = pshEnv.get("_links", {}).get("ssh", {}).get("href")
        sshUrl = sshUrl.replace("ssh://", "")
        # itterate mount points and perform rsync
        appNames = pccProject.getApplicationsParser().getApplicationNames()
        for appName in appNames:
            app = pccProject.getApplication(appName)
            mounts = app.getMounts()
            for _, mountDest in mounts.items():
                self.logger.info("Rsync mount '%s' for application '%s.'" % (mountDest, appName))
                try:
                    app.runCommand(
                        """
                        cd %s
                        chown -R web:web %s
                        """ % (
                            app.APPLICATION_DIRECTORY,
                            mountDest
                        )
                    )
                    app.runCommand(
                        """
                        cd %s
                        rsync -a  -e "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" --max-size=2M %s:%s/ %s/
                        """ % (
                            app.APPLICATION_DIRECTORY,
                            sshUrl,
                            mountDest,
                            mountDest
                        ),
                        "web"
                    )
                except ContainerCommandError:
                    pass

    def syncServices(self, pccProject):
        """ Sync project service assets with Platform.sh. """
        if not self.isRunning():
            raise StateError("Cannot sync services, Platform.sh clone container is not running.") 
        self.logger.info("Fetch service assets.")
        # get project env info
        pshEnv = self._getPshEnv()
        # retrieve ssh url
        sshUrl = pshEnv.get("_links", {}).get("ssh", {}).get("href")
        sshUrl = sshUrl.replace("ssh://", "")
        # itterate relationships and perform asset dumps
        pshRelationsStr = self.runCommand(
            """
            ssh %s -q 'echo $PLATFORM_RELATIONSHIPS | base64 -d 2> /dev/null'
            """ % (
                sshUrl
            )
        )
        pshRelations = json.loads(pshRelationsStr.strip())
        for name in pshRelations:
            self.logger.info("Fetch assets for service relationship '%s.'" % name)
            for relationship in pshRelations[name]:
                fetcher = getPlatformShFetcher(self, pccProject, relationship, sshUrl)
                if not fetcher: continue
                fetcher.fetch()