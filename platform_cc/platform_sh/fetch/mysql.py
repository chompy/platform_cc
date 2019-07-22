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
from .base import PlatformShFetcher
from platform_cc.application.php import PhpApplication
from platform_cc.commands.mysql import getMysqlService

class PlatformShFetchMysql(PlatformShFetcher):
    
    def dump(self):
        return [self._runCommandDump(
            """
            ssh %s -q 'mysqldump -h "%s" -u "%s" --password="%s" %s'
            """ % (
                self.sshUrl,
                self.relationship.get("host", ""),
                self.relationship.get("username", ""),
                self.relationship.get("password", ""),
                self.relationship.get("path", "")
            )
        )]

    def importProject(self, project, dumpPaths = []):
        for dumpPath in dumpPaths:
            service = getMysqlService(project, self.relationship.get("service"))
            if not service: continue
            with open(dumpPath, "rb") as f:
                service.executeSqlDump(
                    self.relationship.get("path", None),
                    f
                )