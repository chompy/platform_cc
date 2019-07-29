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

import docker
import io

class MockDockerVolume:

    def __init__(self, volumes, name, **kwargs):
        self.volumes = volumes
        self.name = name
        self.args = kwargs

    def remove(self):
        self.volumes.volumes.remove(self)

class MockDockerVolumes:

    def __init__(self):
        self.volumes = []

    def create(self, name, **kwargs):
        volume = MockDockerVolume(self, name, **kwargs)
        self.volumes.append(volume)
        return volume

    def get(self, name):
        for i in range(len(self.volumes)):
            if self.volumes[i].name == name:
                return self.volumes[i]
        raise docker.errors.NotFound("TEST")

class MockDockerNetwork:

    def __init__(self, name, **kwargs):
        self.name = name
        self.args = kwargs

class MockDockerNetworks:

    def __init__(self):
        self.networks = []

    def create(self, name, **kwargs):
        network = MockDockerNetwork(name, **kwargs)
        self.networks.append(network)
        return network

    def get(self, name):
        for i in range(len(self.networks)):
            if self.networks[i].name == name:
                return self.networks[i]
        raise docker.errors.NotFound("TEST")

class MockDockerImage:

    def __init__(self, name):
        self.name = name

class MockDockerImages:
    
    def __init__(self):
        self.images = []

    def get(self, name):
        for i in range(len(self.images)):
            if self.images[i].name == name:
                return self.images[i]
        raise docker.errors.ImageNotFound("TEST")

    def pull(self, name):
        image = MockDockerImage(name)
        self.images.append(image)
        return image

    def remove(self, name):
        for i in range(len(self.images)):
            if self.images[i].name == name:
                return self.images.pop(i)

class MockDockerContainer:

    def __init__(self, images, image, **kwargs):
        self.images = images
        self.imageName = image
        self.args = kwargs
        self.stop()

    def start(self):
        self.status = "running"

    def stop(self):
        self.status = "stopped"
    
    def wait(self):
        return

    def remove(self):
        return

    def put_archive(self, path, data):
        if type(data) is not io.BytesIO:
            raise Exception("TEST")

    def commit(self, repo, tag):
        self.images.pull("%s:%s" % (repo, tag))

class MockDockerContainers:

    def __init__(self, images):
        self.images = images
        self.containers = []

    def create(self, image, **kwargs):
        container = MockDockerContainer(
            self.images, image, **kwargs
        )
        self.containers.append(container)
        return container

    def get(self, name):
        for i in range(len(self.containers)):
            if self.containers[i].args.get("name") == name:
                return self.containers[i]
        raise docker.errors.ImageNotFound("TEST")


class MockDocker:

    def __init__(self, timeout=30):
        self.timeout = timeout
        self.images = MockDockerImages()
        self.containers = MockDockerContainers(self.images)
        self.networks = MockDockerNetworks()
        self.volumes = MockDockerVolumes()

    @classmethod
    def from_env(cls, timeout=30):
        return MockDocker(timeout=timeout)