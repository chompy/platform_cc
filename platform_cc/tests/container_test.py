from __future__ import absolute_import
from .base import BaseTest
from ..container import Container
import docker

class ContainerTest(BaseTest):
    """ Test container class. """

    """ Name of test container. """
    CONTAINER_NAME = "_test_container"

    def setUp(self):
        BaseTest.setUp(self)
        self.dockerClient = docker.from_env()
        self.container = Container(
            self.PROJECT_DATA,
            self.CONTAINER_NAME,
            self.dockerClient
        )

    def tearDown(self):
        self.container.getNetwork().remove()
        self.dockerClient.api.close()

    def testContainerName(self):
        """ Ensure expected container name returned. """
        self.assertEqual(
            "%s%s_%s" % (
                Container.CONTAINER_NAME_PREFIX,
                self.PROJECT_DATA["uid"][0:6],
                self.CONTAINER_NAME
            ),
            self.container.getContainerName()
        )

    def testNetworkName(self):
        """ Ensure expected network name returned. """
        self.assertEqual(
            "%s%s" % (
                Container.CONTAINER_NAME_PREFIX,
                self.PROJECT_DATA["uid"][0:6]
            ),
            self.container.getNetworkName()
        )

    def testGetNetwork(self):
        """ Ensure expected Docker network object is returned. """
        network = self.container.getNetwork()
        self.assertIsInstance(
            network,
            docker.models.networks.Network
        )
        self.assertIsNotNone(
            network.id
        )
        self.assertEqual(
            self.container.getNetworkName(),
            network.attrs.get("Name")
        )
        
    def testGetVolume(self):
        """ Test that it's possible to obtain volumes for container. """
        defaultVolume = self.container.getVolume()
        self.assertIsInstance(
            defaultVolume,
            docker.models.volumes.Volume
        )
        self.assertEqual(
            "%s%s_%s" % (
                self.container.CONTAINER_NAME_PREFIX,
                self.PROJECT_DATA["uid"][0:6],
                self.CONTAINER_NAME
            ),
            defaultVolume.name
        )
        defaultVolume.remove()
        namedVolume = self.container.getVolume("test")
        self.assertIsInstance(
            namedVolume,
            docker.models.volumes.Volume
        )
        self.assertEqual(
            "%s%s_%s_test" % (
                self.container.CONTAINER_NAME_PREFIX,
                self.PROJECT_DATA["uid"][0:6],
                self.CONTAINER_NAME
            ),
            namedVolume.name
        )
        namedVolume.remove()

    def testStartContainer(self):
        """ Perform tests against running container. """
        self.assertIsNone(self.container.getContainer())
        self.assertFalse(self.container.isRunning())
        self.container.start()
        container = self.container.getContainer()
        self.assertIsInstance(
            container,
            docker.models.containers.Container
        )
        self.assertTrue(self.container.isRunning())
        self.assertIsNotNone(self.container.getContainerIpAddress())
        self.container.stop()