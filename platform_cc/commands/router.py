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

from __future__ import absolute_import
import os
from cleo import Command
from ..core.router import PlatformRouter

class RouterStart(Command):
    """
    Start the router.

    router:start
    """

    def handle(self):
        router = PlatformRouter()
        router.start()

class RouterStop(Command):
    """
    Stop the router.

    router:stop
    """

    def handle(self):
        router = PlatformRouter()
        router.stop()

class RouterRestart(Command):
    """
    Restart the router.

    router:restart
    """

    def handle(self):
        router = PlatformRouter()
        router.restart()

class RouterAdd(Command):
    """
    Add project to router.

    router:add
    {--p|path=? : Path to project root. (Default=current directory)}
    """

    def handle(self):
        project = getProject(self)
        project.addRouter()

class RouterRemove(Command):
    """
    Remove project from router.

    router:remove
    {--p|path=? : Path to project root. (Default=current directory)}
    {--u|uid=? : Project uid, can be provided instead of path.}
    """

    def handle(self):
        project = getProject(self)
        project.removeRouter()

class RouterNginx(Command):
    """
    Get Nginx configuration for given project.

    router:nginx
    {--p|path=? : Path to project root. (Default=current directory)}
    {--u|uid=? : Project uid, can be provided instead of path.}
    """

    def handle(self):
        project = getProject(self)
        router = project.getRouter()
        router.logger.propagate = False
        nginxConfig = project.buildRouterNginxConfig()
        self.line(nginxConfig)