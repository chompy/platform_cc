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
from .exception.access_token_error import PlatformShAccessTokenError
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

    @classmethod
    def getAccessToken(cls, apiToken):
        """ Retrieve access token from API token. """
        r = requests.post(
            cls.OAUTH_URL,
            json = {
                "client_id" : "platform-api-user",
                "grant_type" : "api_token",
                "api_token" : str(apiToken)
            }
        )
        if not r.text:
            raise PlatformShApiError("API returned empty response.")
        resp = r.json()
        if resp.get("error"):
            raise PlatformShApiError(
                "%s (%s)" % (resp.get("error_description"), resp.get("error"))
            )
        return resp.get("access_token")

    def _apiGet(self, resource):
        """ Perform GET request to given API resource. """
        if not self.config.getAccessToken():
            raise PlatformShAccessTokenError("Access token not set. Please login first with 'platform_sh:login.'")
        r = requests.get(
            "%s/%s" % (self.API_URL, resource),
            headers={"Authorization" : "Bearer %s" % self.config.getAccessToken()}
        )
        return self._handleApiRequest(r)

    def _handleApiRequest(self, r):
        if r.status_code == 401:
            raise PlatformShAccessTokenError("Invalid or expired access token provided. Please refresh your access token with 'platform_sh:login.'")
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
        """ Retrieve or generate SSH key. """
        # use stored ssh key if available
        if self.config.getSshPublicKey() and self.config.getSshPrivateKey():
            return [self.config.getSshPublicKey(), self.config.getSshPrivateKey()]
        # otherwise generate one and upload to platform.sh
        key = rsa.generate_private_key(
            backend=crypto_default_backend(),
            public_exponent=65537,
            key_size=2048
        )
        privateKey = key.private_bytes(
            crypto_serialization.Encoding.PEM,
            crypto_serialization.PrivateFormat.PKCS8,
            crypto_serialization.NoEncryption()
        )
        publicKey = key.public_key().public_bytes(
            crypto_serialization.Encoding.OpenSSH,
            crypto_serialization.PublicFormat.OpenSSH
        )
        uuid = self.getUUID()
        r = requests.post(
            "%s/ssh_keys" % self.API_URL,
            json={
                "value" : publicKey.decode("utf-8"),
                "title" : "%s @ %s" % (self.SSH_KEY_TITLE, platform.uname()[1]),
                "uuid" : uuid
            },
            headers={"Authorization" : "Bearer %s" % self.config.getAccessToken()}
        )
        self._handleApiRequest(r)
        self.config.setSshPublicKey(publicKey.decode("utf-8"))
        self.config.setSshPrivateKey(privateKey.decode("utf-8"))
        return [self.config.getSshPublicKey(), self.config.getSshPrivateKey()]
