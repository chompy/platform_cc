import os
import difflib
import io
import hashlib
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a PHP container. """

    def provision(self):

        # parent method
        DockerProvisionBase.provision(self)

        # install extensions
        print_stdout("  - Install extensions.")
        extensions = self.platformConfig.getRuntime().get("extensions", [])
        extensionConfigs = self.config.get("extensions", {})
        for extensionName in extensions:
            print_stdout("    - %s..." % (extensionName), False)
            extensionConfig = extensionConfigs.get(extensionName, {})
            if not extensionConfig:
                print_stdout("not available.")
                continue
            if extensionConfig.get("core", False):
                print_stdout("already installed (core extension).")
                continue
            depCmdKey = difflib.get_close_matches(
                self.image,
                extensionConfig.keys(),
                1
            )
            if not depCmdKey:
                print_stdout("not available.")
                continue

            self.container.exec_run(
                ["sh", "-c", extensionConfig[depCmdKey[0]]]
            )
            print_stdout("done.")

    def getUid(self):
        """ Generate unique id based on configuration. """
        hashStr = self.image
        hashStr += str(self.platformConfig.getBuildFlavor())
        extensions = self.platformConfig.getRuntime().get("extensions", [])
        extensions.sort()
        extensionConfigs = self.config.get("extensions", {})
        for extension in extensions:
            extensionConfig = extensionConfigs.get(extension, {})
            if not extensionConfig: continue
            if not extensionConfig.get("core", False): continue
            hashStr += extension
        return hashlib.sha256(hashStr).hexdigest()