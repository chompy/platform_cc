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
import tempfile
import hashlib
import base36

class PlatformShFetcher:
    """ Fetch assets required by Platform.sh services. """

    def __init__(self, project, relationship, sshUrl):
        self.project = project
        self.sshUrl = str(sshUrl)
        self.relationship = dict(relationship)
        self.dumpPath = tempfile.gettempdir()

    def _runCommandDump(self, cmd):
        """ Run container command dump the results to temp file. """
        # generate filename for dump file
        dumpFilename = "pcc_%s" % base36.dumps(
            int(
                hashlib.sha256(
                    (
                        "%s%s" % (
                            self.dumpPath,
                            cmd
                        )
                    ).encode("utf-8")
                ).hexdigest(),
                16
            )
        )
        # retrieve docker container
        app = self.project.getApplication()
        dockerContainer = app.getContainer()
        # run command
        (_, output) = dockerContainer.exec_run(
            [
                "sh", "-c", cmd
            ],
            user="web",
            stream=True
        )

        # stream results to file
        dumpFilePath = os.path.join(self.dumpPath, dumpFilename)
        with open(dumpFilePath, "wb") as f:
            res = True
            while (res):
                res = next(output, False)
                if res: f.write(res)
        return dumpFilePath

    def fetch(self):
        """ Fetch asset and import to project. """
        pass

