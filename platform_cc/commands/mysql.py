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

import sys
import os
import io
from cleo import Command
from platform_cc.commands import getProject, outputJson, outputTable
from platform_cc.exception.state_error import StateError

def getMysqlService(project, name = None):
    """ Get MySQL service to use. """
    servicesParser = project.getServicesParser()
    if not name:
        for _name in servicesParser.getServiceNames():
            serviceType = servicesParser.getServiceType(_name).split(":")[0]
            if serviceType not in ["mysql", "mariadb"]: continue
            name = _name
            break
    if not name:
        raise ValueError("No service was specified.")
    serviceType = servicesParser.getServiceType(name).split(":")[0]
    if serviceType not in ["mysql", "mariadb"]:
        raise ValueError(
            "Service '%s' is not a MySQL or MariaDB service." % name
        )
    return project.getService(name)

class MysqlSql(Command):
    """
    Execute SQL commands for MySQL service.

    mysql:sql
        {--p|path=? : Path to project root. (Default=current directory)}
        {--s|service=? : Name of MariaDB service. (Default=first available)}
        {--d|database=?} : Name of database to use.}
    """

    def handle(self):
        project = getProject(self)
        serviceName = self.option("service")
        service = getMysqlService(project, serviceName)
        if not service.isRunning():
            raise StateError(
                "Service '%s' is not running." % serviceName
            )
        cmd = "mysql -h 127.0.0.1 -uroot --password=\"%s\"" % (
            service.getPassword()
        )
        if self.option("database"):
            cmd += " --database=\"%s\"" % self.option("database")
        
        # has stdin
        if not sys.stdin.isatty():
            stdin = sys.stdin
            try:
                stdin = sys.stdin.detach().read()
            except AttributeError:
                pass
            try:
                stdin = stdin.read()
            except AttributeError:
                pass
            byteIo = io.BytesIO(stdin)
            service.uploadFile(
                byteIo,
                "/stdin.txt"
            )
            cmd = ["sh", "-c", "cat /stdin.txt | %s && rm /stdin.txt" % cmd]
            (exitCode, output) = service.getContainer().exec_run(
                cmd,
                user = "root"
            )
            self.line(output.decode("utf-8"))
            return

        service.shell(cmd)

class MysqlDump(Command):
    """
    Execute dump for MySQL service.

    mysql:dump
        {--p|path=? : Path to project root. (Default=current directory)}
        {--s|service=? : Name of MariaDB service. (Default=first available)}
        {database : Name of database to use.}
    """

    def handle(self):
        project = getProject(self)
        serviceName = self.option("service")
        service = getMysqlService(project, serviceName)
        if not service.isRunning():
            raise StateError(
                "Service '%s' is not running." % serviceName
            )
        cmd = "mysqldump -h 127.0.0.1 -uroot --password=\"%s\" %s" % (
            service.getPassword(),
            self.argument("database")
        )
        service.shell(cmd)
