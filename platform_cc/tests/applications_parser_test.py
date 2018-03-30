from __future__ import absolute_import
from .base import BaseTest
from ..parser.applications import ApplicationsParser

class ApplicationsParserTest(BaseTest):
    """
    Test applications parser.
    """

    """ 
    List of expected applications.
    TODO: test multi application setup
    """
    EXPECTED_APPS = ["app"]

    """ List of expected application types. """
    EXPECTED_APP_TYPES = ["php:5.6"]

    def setUp(self):
        self.applicationsParser = ApplicationsParser(
            self.PROJECT_PATH
        )

    def testGetApplicationNames(self):
        """ Test ability to list names of all applications. """
        appNames = self.applicationsParser.getApplicationNames()
        self.assertEqual(
            len(self.EXPECTED_APPS),
            len(appNames)
        )
        for expectedAppName in self.EXPECTED_APPS:
            self.assertIn(
                expectedAppName,
                appNames
            )

    def testGetApplicationConfig(self):
        """ Test ability to get application configuration. """
        appNames = self.applicationsParser.getApplicationNames()
        for appName in appNames:
            appConfig = self.applicationsParser.getApplicationConfiguration(appName)
            self.assertIsInstance(
                appConfig,
                dict
            )
            self.assertIsNotNone(
                appConfig.get("name")
            )
            self.assertIsNotNone(
                appConfig.get("type")
            )

if __name__ == "__main__":
    unittest.main()