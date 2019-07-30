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
from platform_cc.variables.var_dict import DictVariables
from platform_cc.variables.var_json_file import JsonFileVariables
from platform_cc.variables import getVariableStorage

class TestVar(unittest.TestCase):

    """ Test var classes. """

    def testVarDict(self):
        var = DictVariables(config={
            "dict_vars" : {
                "a" : 1,
                "b" : 2,
                "c" : 3,
            }
        })
        # test 'get'
        self.assertEqual(var.get("a"), 1)
        self.assertEqual(var.get("b"), 2)
        self.assertEqual(var.get("c"), 3)
        self.assertEqual(var.get("d"), None)
        self.assertEqual(var.get("d", 4), 4)
        # test 'all'
        allVar = var.all()
        self.assertEqual(allVar.get("a"), 1)
        # ensure 'allVar' is a copy and does not affect the original var obj
        allVar["a"] = 2
        self.assertEqual(var.get("a"), 1)
        self.assertEqual(allVar.get("a"), 2)
        # test 'delete'
        var.delete("a")
        self.assertEqual(var.get("a"), None)

    def testGlobalVarDict(self):
        var = DictVariables(config={
            "dict_vars" : {
                "a" : 1,
                "b" : 2,
                "c" : 3,
            },
            "global_vars" : {
                "a" : 10,
                "d" : 4,
            }
        })
        self.assertEqual(var.get("a"), 1, "ensure global var doesn't overwrite local var")
        self.assertEqual(var.get("d"), 4, "ensure get falls back to global var if local doesn't exist")
        var.delete("a")
        self.assertEqual(var.get("a"), 10, "ensure fallback to global var if local var deleted")
        var.delete("d")
        self.assertEqual(var.get("d"), 4, "ensure delete doesn't delete global var")

    def testJsonVarDict(self):
        # value error raised if json_path not passed
        with self.assertRaises(ValueError):
            JsonFileVariables()
        # check that vars are fetchable
        var = JsonFileVariables({
            "json_path" : os.path.join(os.path.dirname(__file__), "data", JsonFileVariables.JSON_STORAGE_FILENAME)
        })
        self.assertEqual(var.get("a"), 1)
        self.assertEqual(var.get("b"), 2)
        self.assertEqual(var.get("c"), 3)
        # check that var is settable
        currentTime = time.time()
        var.set("time", currentTime)
        self.assertEqual(var.get("time"), currentTime, "set time in json var and ensure fetch same value")
        var2 = JsonFileVariables({
            "json_path" : os.path.join(os.path.dirname(__file__), "data", JsonFileVariables.JSON_STORAGE_FILENAME)
        })
        self.assertEqual(var2.get("time"), currentTime, "ensure time json var is fetchable after reloading object")

    def testGetVar(self):
        self.assertIsInstance(getVariableStorage(config={"storage_handler": "dict"}), DictVariables)
        self.assertIsInstance(getVariableStorage(config={"storage_handler": "json", "json_path" : "./test.txt"}), DictVariables)
        with self.assertRaises(NotImplementedError):
            getVariableStorage(config={"storage_handler" : "unknown"})

if __name__ == '__main__':
    unittest.main()