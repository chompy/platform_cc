import os
from provision_base import DockerProvisionBase

class DockerProvision(DockerProvisionBase):

    """ Provision a PHP container. """

    def provision(self):
        
        # install ssh key
        print "  - Install SSH key file...",
        self.container.exec_run(
            ["mkdir", "/creds"]
        )
        self.copyFile(
            os.path.join(
                self.platformConfig.getDataPath(),
                ".id_rsa"
            ),
            "/creds/id_rsa"
        )
        print "done."

        # add 'web' user
        print "  - Create 'web' user...",
        password = self.randomString(10)
        self.container.exec_run(
            ["useradd", "-d", "/app", "-m", "-p", password, "web"]
        )
        print "done."

        # rsync app
        print "  - Copy application to container...",
        self.container.exec_run(
            ["rsync", "-a", "--exclude", ".platform", "--exclude", ".git", "--exclude", ".platform.app.yaml", "/mnt/app/", "/app"]
        )
        print "done."
