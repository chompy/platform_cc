import os
from provision_base import DockerProvisionBase

class DockerProvision(DockerProvisionBase):

    """ Provision a PHP container. """

    def provision(self):
        # install ssh key
        print self.container.exec_run(
            ["mkdir", "/creds"]
        )
        self.copyFile(
            os.path.join(
                self.platformConfig.getDataPath(),
                ".id_rsa"
            ),
            "/creds/id_rsa"
        )
        # install git
        password = self.randomString(10)
        print self.container.exec_run(
            ["apt-get", "update"]
        )
        print self.container.exec_run(
            ["apt-get", "install", "-y", "git"]
        )
        print self.container.exec_run(
            ["apt-get", "clean"]
        )
        # add 'web' user
        print self.container.exec_run(
            ["useradd", "-d", "/app", "-m", "-p", password, "web"]
        )