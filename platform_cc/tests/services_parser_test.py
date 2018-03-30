from __future__ import absolute_import
from .base import BaseTest
from ..parser.services import ServicesParser

class ServicesParserTest(BaseTest):
    """
    Test services parser.
    """

    """ List of expected services. """
    EXPECTED_SERVICES = ["mysqldb", "rediscache", "redisdata", "memcached"]

    """ List of expected service types. """
    EXPECTED_SERVICE_TYPES = ["mysql:10.1", "redis:3.2", "redis-persistent:3.2", "memcached:1.4"]

    def setUp(self):
        self.servicesParser = ServicesParser(self.PROJECT_PATH)

    def testGetServiceNames(self):
        """ Test ability to list names of all services. """
        serviceNames = self.servicesParser.getServiceNames()
        self.assertEqual(
            len(self.EXPECTED_SERVICES),
            len(serviceNames)
        )
        for expectedServiceName in self.EXPECTED_SERVICES:
            self.assertIn(
                expectedServiceName,
                serviceNames
            )

    def testGetServiceType(self):
        """ Test that expected service type is returned. """
        for i in range(len(self.EXPECTED_SERVICES)):
            self.assertEqual(
                self.EXPECTED_SERVICE_TYPES[i],
                self.servicesParser.getServiceType(
                    self.EXPECTED_SERVICES[i]
                )
            )

    def testGetServiceConfiguration(self):
        """ That that expected service configuration is returned. """
        for serviceName in self.EXPECTED_SERVICES:
            serviceConfig = self.servicesParser.getServiceConfiguration(serviceName)
            self.assertIsInstance(
                serviceConfig,
                dict
            )
            self.assertEqual(
                serviceConfig.get("_name"),
                serviceName
            )
            self.assertIsNotNone(
                serviceConfig.get("_type")
            )

if __name__ == "__main__":
    unittest.main()