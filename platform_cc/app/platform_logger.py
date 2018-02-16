from cleo import Output

class PlatformLogger:

    """ Handles logging to stdout. (And maybe eventually a file.) """

    def __init__(self, command = None):
        self.command = command

    def logEvent(self, message, indent = 0, verbosity = Output.VERBOSITY_NORMAL):
        """ Log a standard event. (App start up, commmand ran, etc) """
        if verbosity > self.command.output.get_verbosity(): return
        if indent < 0 or not message: return
        indentStr = ">"
        if indent > 0:
            indentStr = "-".rjust(indent * 3)
        self.command.line(
            "%s %s" % (
                indentStr,
                message
            )
        )

    def printContainerOutput(self, results, verbosity = Output.VERBOSITY_VERBOSE):
        """ Log the output of a command from a Docker container. """
        if verbosity > self.command.output.get_verbosity(): return
        self.command.line(
            "= COMMAND OUTPUT =========================================\n%s\n========================================================" % (
                str(results)
            )
        )