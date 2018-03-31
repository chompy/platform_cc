from __future__ import absolute_import
import os
import sys
import unittest
from unittest import TestCase
sys.path.append(
    os.path.join(
        os.path.dirname(__file__),
        ".."
    )
)

class BaseTest(TestCase):
    """
    Base class for all test.
    """

    """ Path to 'sample' project. """
    PROJECT_PATH = os.path.join(
        os.path.dirname(__file__),
        "sample_project"
    )

    """ Sample project data. """
    PROJECT_DATA = {
        "path"      : PROJECT_PATH,
        "uid"       : "123abc456def",
        "entropy"   : "123random"
    }
