import docker
import dockerpty

class BasePlatformService:
    """
    Base class for Platform.sh services.
    """
    
    def _buildContainer(self):
        pass