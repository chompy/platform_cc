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
from platform_cc.parser.routes import RoutesParser
from platform_cc.parser.services import ServicesParser
from platform_cc.exception.parser_error import ParserError

class TestParser(unittest.TestCase):

    """ Test parser classes. """

    TEST_PROJECT_UID = "abcdefg1234567"

    def _getApplicationParser(self):
        return ApplicationsParser(
            os.path.join(
                os.path.dirname(__file__),
                "data"
            )
        )

    def _getRoutesParser(self):
        return RoutesParser({
            "path" : os.path.join(
                    os.path.dirname(__file__),
                    "data"
                ),
            "uid" : self.TEST_PROJECT_UID,
            "short_uid" : self.TEST_PROJECT_UID[:5]
        })

    def _getServicesParser(self):
        return ServicesParser(
            os.path.join(
                os.path.dirname(__file__),
                "data"
            )
        )

    def testAppParser(self):
        app = self._getApplicationParser()
        self.assertIn(
            "test_app", app.getApplicationNames()
        )
        with self.assertRaises(ParserError):
            app.getApplicationConfiguration("test_app_none")
        config = app.getApplicationConfiguration("test_app")
        self.assertEqual(config.get("variables").get("env").get("TEST_ENV"), "yes")
        self.assertEqual(config.get("relationships").get("database"), "mysqldb:mysql")
        self.assertEqual(config.get("web").get("locations").get("/").get("passthru"), "/index.php")
        # test pcc override file
        self.assertEqual(config.get("variables").get("env").get("IS_PCC"), "yes")

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

    def testRouteParser(self):
        routeParser = self._getRoutesParser()
        routes = routeParser.getRoutes()
        self.assertIn("www.example.com", routes[0].get("hostnames"))
        self.assertIn("www.example.com.%s.*" % self.TEST_PROJECT_UID[:5], routes[0].get("hostnames"))
        self.assertEqual("upstream", routes[0].get("type"))
        self.assertEqual("test_app", routes[0].get("upstream"))
        self.assertTrue(routes[0].get("cache").get("enabled"))        
        self.assertEqual(routes[0].get("redirects").get("paths").get("/test2").get("to"), "/test1")
        self.assertEqual("redirect", routes[1].get("type"))
        self.assertEqual("https://www.example.com", routes[1].get("to"))
        routesByHostname = routeParser.getRoutesByHostname()
        self.assertEqual(1, len(routesByHostname.get("www.example.com")))
        self.assertEqual(1, len(routesByHostname.get("example.com")))
        self.assertEqual("upstream", routesByHostname.get("www.example.com")[0].get("type"))
        self.assertEqual("redirect", routesByHostname.get("example.com")[0].get("type"))
        routesEnv = routeParser.getRoutesEnvironmentVariable()
        self.assertIn("https://www.example.com/", routesEnv)
        self.assertEqual("upstream", routesEnv.get("https://www.example.com/").get("type"))
        self.assertIn("https://www.example.com.%s.*/" % self.TEST_PROJECT_UID[:5], routesEnv)
        self.assertEqual("upstream", routesEnv.get("https://www.example.com.%s.*/" % self.TEST_PROJECT_UID[:5]).get("type"))
        self.assertIn("https://example.com/", routesEnv)
        self.assertEqual("redirect", routesEnv.get("https://example.com/").get("type"))
        self.assertIn("https://example.com.%s.*/" % self.TEST_PROJECT_UID[:5], routesEnv)
        self.assertEqual("redirect", routesEnv.get("https://example.com.%s.*/" % self.TEST_PROJECT_UID[:5]).get("type"))
        # test pcc override file
        self.assertEqual(1, len(routesByHostname.get("test.example.com")))
        self.assertEqual("redirect", routesByHostname.get("test.example.com")[0].get("type"))

    def testServicesParser(self):
        servicesParser = self._getServicesParser()
        serviceNames = servicesParser.getServiceNames()
        self.assertIn("mysqldb", serviceNames)
        self.assertIn("memcached", serviceNames)
        self.assertEqual("mysql:10.0", servicesParser.getServiceType("mysqldb"))
        self.assertEqual("memcached:1.4", servicesParser.getServiceType("memcached"))
        mysqlConf = servicesParser.getServiceConfiguration("mysqldb")
        self.assertIn("main", mysqlConf.get("schemas"))
        self.assertEqual("main", mysqlConf.get("endpoints").get("mysql").get("default_schema"))
        # test pcc override file
        self.assertIn("test", mysqlConf.get("schemas"))
        self.assertEqual("test", mysqlConf.get("endpoints").get("test_endpoint").get("default_schema"))

if __name__ == '__main__':
    unittest.main()