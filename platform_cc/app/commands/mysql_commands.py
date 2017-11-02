import os
from cleo import Command
from app.commands import getProject, getAppsToInvoke

def findMysqlService(command):
    appName = command.option("app")
    serviceName = command.option("service")
    project = getProject(command, False)
    application = None
    for app in project.getApplications():
        if app.config.getName() == appName or not appName:
            application = app
            break
    if not application:
        print "ERROR: Application '%s' does not exist." % appName
        return None
    service = None
    for sv in application.getServices():
        if sv.config.getType()[:5] == "mysql" and (sv.config.getName() == serviceName or not serviceName):
            service = sv
    if not service:
        print "ERROR: MySQL service '%s' does not exist." % serviceName
    return service

class MysqlSql(Command):
    """
    Run SQL on MySQL database

    mysql:sql
        {--a|app=? : Application where database service resides. (First available if not provided.)}
        {--s|service=? : Database service. (First available if not provided.)}
        {--d|database=? : Database to select. }
        {--dp|dumppath=? : Path to SQL dump to import.}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        service = findMysqlService(self)
        if not service: return
        dumpPath = self.option("dumppath")
        database = self.option("database")
        password = service.docker.getProvisioner().getPassword()
        if dumpPath:
            if not os.path.exists(dumpPath):
                print "ERROR: SQL dump at '%s' does not exists." % dumpPath
                return
            service.docker.getProvisioner().copyFile(
                dumpPath,
                "/dump.sql"
            )
            service.shell(
                "bash -c 'mysql -uroot --password=\"%s\"%s < /dump.sql'" % (
                    password,
                    ((" %s" % database) if database else "")
                )
            )
            return
        service.shell(
            "mysql -uroot --password=\"%s\"%s" % (
                password,
                ((" %s" % database) if database else "")
            )
        )

class MysqlDump(Command):

    """
    Dump database in MySQL database

    mysql:dump
        {database : Database to dump.}
        {--a|app=? : Application where database service resides. (First available if not provided.)}
        {--s|service=? : Database service. (First available if not provided.)}
        {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        service = findMysqlService(self)
        if not service: return
        password = service.docker.getProvisioner().getPassword()
        database = self.argument("database")
        service.shell("mysqldump -uroot --password=\"%s\" %s" % (
            database,
            password
        ))