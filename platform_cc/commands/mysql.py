import sys
import os
import io
from cleo import Command
from platform_cc.commands import getProject, outputJson, outputTable
from platform_cc.exception.state_error import StateError

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
        servicesParser = project.getServicesParser()
        if not serviceName:
            for _serviceName in servicesParser.getServiceNames():
                serviceType = servicesParser.getServiceType(_serviceName).split(":")[0]
                if serviceType not in ["mysql", "mariadb"]: continue
                serviceName = _serviceName
                break
        if not serviceName:
            raise ValueError("No service was specified.")
        serviceType = servicesParser.getServiceType(serviceName).split(":")[0]
        if serviceType not in ["mysql", "mariadb"]:
            raise ValueError(
                "Service '%s' is not a MySQL or MariaDB service." % serviceName
            )
        service = project.getService(serviceName)
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
