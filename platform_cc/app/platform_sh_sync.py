
class PlatformShSync:

    def __init__(self, application, sshCmd, logger=None):
        self.application = application
        self.sshUrl = sshCmd.split(" ")[-1]
        self.logger = logger
        self.application.copySshKey()

    def _runCmd(self, cmd):
        """ Run shell command on Platform.sh """
        return self.application.docker.getContainer().exec_run(
            ["ssh", "-o", "StrictHostKeyChecking=no", self.sshUrl, cmd],
            user="web",
            privileged=True
        )

    def rsyncCopy(self, path, limitFileTypes = []):
        """ Perform a rsync from Platform.sh. """

        # limit to files with provided extensions
        includeCmd = []
        if limitFileTypes:
            includeCmd.append('--include="*/"')
            for value in limitFileTypes:
                includeCmd.append('--include="*.%s" ' % value)
            includeCmd.append('--exclude="*"')
        # run cmd
        results = self.application.docker.getContainer().exec_run(
            [
                "rsync",
                "-zar",
                "%s:%s" % (self.sshUrl, path),
                path
            ] + includeCmd,
            user="web",
            privileged=True
        )
        if results and self.logger:
            self.logger.printContainerOutput(
                results
            )