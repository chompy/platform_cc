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

import unittest
import io
import docker
import logging
from .mock_docker import MockDocker
from platform_cc.core.container import Container
from platform_cc.exception.state_error import StateError
from platform_cc.exception.container_command_error import ContainerCommandError

class TestContainer(unittest.TestCase):

    """ Test container class. """

    def _createMockContainer(self):
        class MockVolumeContainer(Container):
            def getContainerVolumes(self):
                return {
                    self.getVolumeName(): {
                        "bind": "/var/lib/vol1",
                        "mode": "rw"
                    }
                }
        return MockVolumeContainer(
            {
                "uid" : "abcdefgh",
                "short_uid" : "abcd"
            },
            "test",
            MockDocker()
        )

    def testContainerStart(self):
        container = self._createMockContainer()
        container.start()
        self.assertEqual(
            container.getContainer().status,
            "running",
            "container should be running"
        )
        self.assertTrue(
            container.isRunning()
        )

    def testGetNetwork(self):
        container = self._createMockContainer()
        network = container.getNetwork()
        self.assertEqual(
            network.name,
            "%s%s" % (Container.CONTAINER_NAME_PREFIX, container.project.get("short_uid")),
            "network has expected name"
        )
        self.assertTrue(
            Container.LABEL_PREFIX in network.args["labels"].keys(),
            "network has PCC label"
        )

    def testUploadFile(self):
        container = self._createMockContainer()
        fileObj = io.BytesIO(bytes("TEST".encode("utf-8")))
        self.assertEqual(fileObj.tell(), 0, "upload file object not yet read")
        # should fail if container isn't started
        with self.assertRaises(StateError):
            container.uploadFile(fileObj, "/app/test")
        self.assertEqual(fileObj.tell(), 0, "upload file object not read when container isn't running")
        container.start()
        container.uploadFile(fileObj, "/app/test")
        self.assertEqual(fileObj.tell(), 4, "upload file object has been read")
        
    def testGetImage(self):
        container = self._createMockContainer()
        with self.assertRaises(docker.errors.NotFound):
            container.docker.images.get(container.getBaseImage())
        container.pullImage()
        image = container.docker.images.get(container.getBaseImage())
        self.assertEqual(
            image.name,
            container.getBaseImage(),
            "container pull image pulls the base image"
        )

    def testCommit(self):
        container = self._createMockContainer()
        container.start()
        container.commit()
        image = container.docker.images.get(container.getDockerImage())
        self.assertEqual(
            image.name, container.getDockerImage()
        )

    def testPurge(self):
        container = self._createMockContainer()
        container.start()
        container.commit()
        container.docker.images.get(container.getDockerImage())
        container.docker.volumes.get(container.getVolumeName())
        container.purge()
        with self.assertRaises(docker.errors.NotFound):
            container.docker.images.get(container.getDockerImage())
        with self.assertRaises(docker.errors.NotFound):
            container.docker.volumes.get(container.getVolumeName())

    def testVolumes(self):
        class MockVolumeContainer(Container):
            def getContainerVolumes(self):
                return {
                    self.getVolumeName(): {
                        "bind": "/var/lib/vol1",
                        "mode": "rw"
                    },
                    "test_vol2": {
                        "bind": "/var/lib/vol1",
                        "mode": "rw"
                    }
                }
        container = MockVolumeContainer(
            {
                "uid" : "abcdefgh",
                "short_uid" : "abcd"
            },
            "test",
            MockDocker()
        )
        for key in container.getContainerVolumes():
            with self.assertRaises(docker.errors.NotFound):
                container.docker.volumes.get(key)
        container.start()
        for key in container.getContainerVolumes():
            vol = container.docker.volumes.get(key)
            self.assertEqual(vol.name, key)
        container.stop()
        container.purge()
        for key in container.getContainerVolumes():
            with self.assertRaises(docker.errors.NotFound):
                container.docker.volumes.get(key)

    def testRunCommand(self):
        container = self._createMockContainer()
        logging.disable(logging.ERROR)
        with self.assertRaises(StateError):
            container.runCommand("test cmd")
        container.start()
        output = container.runCommand("test cmd")
        self.assertEqual(output, "PASS", "run command output PASS")
        with self.assertRaises(ContainerCommandError):
            output = container.runCommand("test cmd", user="error")

    def testShell(self):
        container = self._createMockContainer()
        with self.assertRaises(StateError):
            container.shell()
        # can't shell in to mock docker but ensure an attempt is made
        # and an expected exception is thrown
        container.start()
        with self.assertRaises(io.UnsupportedOperation):
            container.shell()
        # test stdin
        dataStr = "TEST-123"
        data = io.BytesIO(dataStr.encode("utf-8"))
        self.assertEqual(data.tell(), 0, "new stdin data at position 0")
        container.shell(stdin=data)
        self.assertEqual(data.tell(), len(dataStr), "stdin passed in to shell has been read")

if __name__ == '__main__':
    unittest.main()