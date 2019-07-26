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

from .exception.api_error import PlatformShApiError
from .exception.access_error import PlatformShAccessError
from .exception.config_error import PlatformShConfigError
from .config import PlatformShConfig
from cryptography.hazmat.primitives import serialization as crypto_serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.backends import default_backend as crypto_default_backend
import os
import json
import requests
import platform

class PlatformShApi:
    """
    Interact with Platform.sh's API.
    """

    API_URL = "https://api.platform.sh"
    OAUTH_URL = "https://accounts.platform.sh/oauth2/token"
    SSH_KEY_TITLE = "Platform.CC"

    def __init__(self, config = None):
        if not config:
            config = PlatformShConfig()
        self.config = config
        self.uuid = ""

    def _fetchAccessToken(self):
        """ Fetch access token and update config. """
        # must have API token
        if not self.config.getApiToken():
            raise PlatformShAccessError("API token not set. Please login first with 'platform_sh:login.'")
        r = requests.post(
            self.OAUTH_URL,
            json = {
                "client_id" : "platform-api-user",
                "grant_type" : "api_token",
                "api_token" : self.config.getApiToken()
            }
        )
        if not r.text:
            raise PlatformShApiError("API returned empty response while trying to retrieve access token.")
        resp = r.json()
        if resp.get("error"):
            raise PlatformShApiError(
                "%s (%s)" % (resp.get("error_description"), resp.get("error"))
            )
        self.config.setAccessToken(resp.get("access_token"))

    def _apiGet(self, resource, tryAccessTokenFetch=True):
        """ Perform GET request to given API resource. """
        if not self.config.getAccessToken():
            self._fetchAccessToken()
        r = requests.get(
            "%s/%s" % (self.API_URL, resource),
            headers={"Authorization" : "Bearer %s" % self.config.getAccessToken()}
        )
        if r.status_code == 401 and tryAccessTokenFetch:               
            self._fetchAccessToken()
            return self._apiGet(resource, tryAccessTokenFetch=False)
        return self._handleApiRequest(r)

    def _handleApiRequest(self, r):
        """ Handle an API request. """
        if r.status_code == 401:
            raise PlatformShAccessError("Invalid or expired access token provided. Please refresh your API token with 'platform_sh:login.'")
        r.raise_for_status()
        if not r.text:
            raise PlatformShApiError("API returned empty response.")
        resp = r.json()
        if type(resp) is dict and resp.get("error"):
            raise PlatformShApiError(
                "%s (%s)" % (resp.get("error_description"), resp.get("error"))
            )
        return resp

    def getUUID(self):
        """ Retrieve user UUID. """
        if self.uuid: return self.uuid
        resp = self._apiGet("me")
        self.uuid = resp.get("uuid")
        return self.uuid

    def getProjects(self):
        """ Retrieve list of projects. """
        return self._apiGet("projects")

    def getProjectInfo(self, projectId):
        """ Retrieve project info. """
        return self._apiGet(
            "projects/%s" % str(projectId)
        )

    def getEnvironmentInfo(self, projectId, environmentId = "master"):
        """ Retrieve environment info. """
        return self._apiGet(
            "projects/%s/environments/%s" % (
                projectId,
                environmentId
            )
        )

    def getDeployment(self, projectId, environmentId = "master"):
        """ Retrieve deployment details. """
        return self._apiGet(
            "projects/%s/environments/%s/deployments" % (
                str(projectId),
                str(environmentId)
            )
        )

    def getSshKeypair(self):
        """ Retrieve SSH key. """
        if not self.config.getSshPrivateKey() or not self.config.getSshPublicKey():
            raise PlatformShConfigError("Missing SSH private or public key. Please set your SSH private key with 'platform_sh:set_ssh.'")
        return [self.config.getSshPublicKey(), self.config.getSshPrivateKey()]
