from __future__ import absolute_import
import os
from .base import BaseTest
from ..variables.var_json import JsonVariables

class VariablesTest(BaseTest):
    """
    Test variables handlers.
    """

    def setUp(self):
        BaseTest.setUp(self)
        self.jsonVariables = JsonVariables(
            self.PROJECT_PATH,
            self.PROJECT_DATA
        )

    def tearDown(self):
        os.remove(
            os.path.join(
                self.PROJECT_PATH,
                self.jsonVariables.JSON_STORAGE_FILENAME
            )
        )

    def testVariableSet(self):
        """ Test to ensure variables can be set correctly. """
        self.jsonVariables.set("test", "1")
        self.jsonVariables.set("test2", "hello")
        self.assertEqual(
            "1",
            self.jsonVariables.get("test")
        )
        self.assertEqual(
            "hello",
            self.jsonVariables.get("test2")
        )

    def testVariableList(self):
        """ Test to ensure variable list can be obtained correctly. """
        self.jsonVariables.set("test", "1")
        self.jsonVariables.set("test2", "hello")
        self.jsonVariables.set("test3", "abcdef")
        allVars = self.jsonVariables.all()
        self.assertEqual(
            3,
            len(allVars)
        )
        self.assertEqual(
            "1",
            allVars.get("test")
        )
        self.assertEqual(
            "hello",
            allVars.get("test2"),
        )
        self.assertEqual(
            "abcdef",
            allVars.get("test3")
        )