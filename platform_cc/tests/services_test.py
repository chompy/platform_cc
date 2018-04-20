from __future__ import absolute_import
from .base import BaseTest
from ..parser.services import ServicesParser
from ..services import getService
from ..services.memcached import MemcachedService
from ..services.mariadb import MariaDbService

class ServicesTest(BaseTest):
    """
    Test service handlers.
    """

    def setUp(self):
        BaseTest.setUp(self)
        self.servicesParser = ServicesParser(self.PROJECT_PATH)

    def testMemcachedService(self):
        """ Test memcached service. """
        memcachedService = getService(
            self.PROJECT_DATA,
            self.servicesParser.getServiceConfiguration("memcached")
        )
        self.assertIsInstance(
            memcachedService,
            MemcachedService
        )
        self.assertEqual(
            "memcached:1",
            memcachedService.getDockerImage()
        )
        self.assertEqual(
            "memcached:1.4",
            memcachedService.getType()
        )

    def testMariaDbService(self):
        """ Test mariadb service. """
        mariaDbService = getService(
            self.PROJECT_DATA,
            self.servicesParser.getServiceConfiguration("mysqldb")
        )
        self.assertEqual(
            "mariadb:10.1",
            mariaDbService.getDockerImage()
        )
        self.assertEqual(
            "mysql:10.1",
            mariaDbService.getType()
        )
        self.assertIsInstance(
            mariaDbService.getPassword(),
            str
        )
        self.assertEqual(
            mariaDbService.getContainerEnvironmentVariables().get("MYSQL_ROOT_PASSWORD"),
            mariaDbService.getPassword()
        )
        self.assertNotEqual(
            mariaDbService.getPassword(),
            mariaDbService.getPassword("test")
        )
        #TODO volumes, container start?