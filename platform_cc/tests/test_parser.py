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

import os
import unittest
import time
from platform_cc.parser.applications import ApplicationsParser
from platform_cc.exception.parser_error import ParserError

class TestParser(unittest.TestCase):

    """ Test parser classes. """

    def _getApplicationParser(self):
        return ApplicationsParser(
            os.path.join(
                os.path.dirname(__file__),
                "data"
            )
        )

    def testAppName(self):
        app = self._getApplicationParser()
        self.assertIn(
            "test_app", app.getApplicationNames()
        )

    def testAppConfig(self):
        app = self._getApplicationParser()
        with self.assertRaises(ParserError):
            app.getApplicationConfiguration("test_app_none")
        config = app.getApplicationConfiguration("test_app")
        self.assertEqual(config.get("variables").get("env").get("TEST_ENV"), "yes")
        self.assertEqual(config.get("relationships").get("database"), "mysqldb:mysql")
        self.assertEqual(config.get("web").get("locations").get("/").get("passthru"), "/index.php")

    def testMultiApp(self):
        app = ApplicationsParser(
            os.path.join(
                os.path.dirname(__file__),
                "data",
                "multi-app"
            )
        )
        self.assertEqual(len(app.getApplicationNames()), 2)
        self.assertIn("test_app1", app.getApplicationNames())
        self.assertIn("test_app2", app.getApplicationNames())
        with self.assertRaises(ParserError):
            app.getApplicationConfiguration("test_app")
        config1 = app.getApplicationConfiguration("test_app1")
        config2 = app.getApplicationConfiguration("test_app2")
        self.assertEqual(config1.get("name"), "test_app1")
        self.assertEqual(config2.get("name"), "test_app2")


if __name__ == '__main__':
    unittest.main()