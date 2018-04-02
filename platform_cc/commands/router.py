import os
from cleo import Command
from router import PlatformRouter

class RouterStart(Command):
    """
    Start the router.

    router:start
    """

    def handle(self):
        router = PlatformRouter()
        router.start()
        self.line(router.getContainerName())

class RouterStop(Command):
    """
    Stop the router.

    router:stop
    """

    def handle(self):
        router = PlatformRouter()
        router.stop()
        self.line(router.getContainerName())
