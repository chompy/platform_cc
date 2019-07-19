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
from platform_cc.container import Container
from platform_cc.application.php import PhpApplication
from .config import PlatformShConfig
from .api import PlatformShApi
from .fetch import getPlatformShFetcher
from platform_cc.project import PlatformProject
from .exception.api_error import PlatformShApiError

class PlatformShCloner(Container):

    """ Tools to clone Platform.sh environment. """

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

    def getBaseImage(self):
        # just use php as the image contains ssh/
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

    def clone(self):
        self.logger.info("Retrieve project information (%s)." % self.project.get("_project_id"))
        pshApi = PlatformShApi(self.config)
        # get project info
        pshProject = pshApi.getProjectInfo(self.project.get("_project_id"))
        # retrieve git url
        pshGitUrl = pshProject.get("repository", {}).get("url", "")
        if not pshGitUrl:
            raise PlatformShApiError("Could not fetch Git URL for project '%s'" % self.project.get("_project_id"))
        parsedGitUrl = urllib.parse.urlparse("ssh://%s" % pshGitUrl)
        # get project env info
        pshEnv = pshApi.getEnvironmentInfo(self.project.get("_project_id"), self.project.get("_environment"))
        # retrieve ssh url
        sshUrl = pshEnv.get("_links", {}).get("ssh", {}).get("href")
        parsedSshUrl = urllib.parse.urlparse(sshUrl)
        sshUrl = sshUrl.replace("ssh://", "")
        # start clone container
        self.start()
        # setup ssh
        cmd = """
        mkdir -p ~/.ssh
        echo "$PSH_SSH_KEY" > ~/.ssh/id_rsa
        chmod 0600 ~/.ssh/id_rsa
        ssh-keyscan %s > ~/.ssh/known_hosts
        ssh-keyscan %s >> ~/.ssh/known_hosts
        chmod 0600 ~/.ssh/known_hosts
        """ % (
            parsedSshUrl.hostname,
            parsedGitUrl.hostname
        )
        try:
            self.runCommand(cmd)
        except Exception as e:
            self.stop()
            raise e

        # git clone
        if not os.path.exists(os.path.join(self.project.get("_path"), self.project.get("_project_id"))):
            self.logger.info(
                "Cloning project %s (%s)." % (
                    self.project.get("_project_id"),
                    self.project.get("_environment")
                )
            )
            cmd = """
            mkdir -p %s
            cd %s
            git clone "%s" --branch "%s"
            """ % (
                PhpApplication.APPLICATION_DIRECTORY,
                PhpApplication.APPLICATION_DIRECTORY,
                pshGitUrl,
                self.project.get("_environment")
            )
            try:
                self.runCommand(cmd)
            except Exception as e:
                self.stop()
                raise e

        # start project
        self.logger.info("Start project.")
        project = PlatformProject.fromPath(
            os.path.join(
                self.project.get("_path"),
                self.project.get("_project_id")
            )
        )
        project.start()

        # ssh access to fetch assets
        self.logger.info("Fetch and import assets.")
        cmd = """
        ssh %s -q 'echo $PLATFORM_RELATIONSHIPS | base64 -d 2> /dev/null'
        """ % (
            sshUrl
        )
        try:
            pshRelationsStr = self.runCommand(cmd)
            pshRelations = json.loads(pshRelationsStr.strip())
            for name in pshRelations:
                for relationship in pshRelations[name]:
                    fetcher = getPlatformShFetcher(relationship, sshUrl)
                    if not fetcher: continue
                    cmd = fetcher.getFetchCommand()
                    if cmd:
                        self.runCommand(cmd)
                        

        except Exception as e:
            self.stop()
            raise e

        self.stop()