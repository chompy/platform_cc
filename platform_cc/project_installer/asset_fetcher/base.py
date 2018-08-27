import os
import tempfile

class BaseAssetFetcher:
    """
    Base asset fetcher class.
    """

    def __init__(self, config):
        self.config = config
        self.path = None

    def __del__(self):
        if self.path:
            os.remove(self.path)

    def get(self):
        """ Fetch asset and store it in temporary file. """
        if self.path: return self.path
        return None