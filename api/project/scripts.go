/*
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
*/

package project

// appContainerCmd is the application container start command.
const appContainerCmd = `
usermod -u %d app
groupmod -g %d app
usermod -u %d web
groupmod -g %d web
umount /etc/hosts
umount /etc/resolv.conf
mkdir -p /run/shared /run/rpc_pipefs/nfs
cat >/tmp/fake-rpc.py <<EOF
from gevent.monkey import patch_all;
patch_all();
from gevent_jsonrpc import RpcServer;
import json;
def rootFactory(c, a):
	with open("/tmp/.ready", "w") as f: f.write("true")
	c.send(json.dumps({"jsonrpc":"2.0","result":True,"id": json.loads(c.recv(1024))["id"]}))
RpcServer(
	"/run/shared/agent.sock",
	"foo",
	root=None,
	root_factory=rootFactory
)._accepter_greenlet.get();
EOF
python /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
runsvdir -P /etc/service &> /tmp/runsvdir.log &
chown -R web:web /run
until [ -f /run/config.json ]; do sleep 1; done
/etc/platform/boot
exec init
`

// appBuildCmd is the build command for applications.
const appBuildCmd = `
cat >/tmp/build.py <<EOF
from platformsh_gevent import patch ; patch()
import os
import sys
import json
kwargs = json.load(sys.stdin)
from platformsh_app.builder.log import ActivityLogger
from platformsh_app.builder import Builder
from platformsh_app.builder import build
log = ActivityLogger(sys.stdout)
builder = Builder.from_application(
    log=log,
    **kwargs
)
builder._generate_configuration()
builder._drop_privileges()
os.chdir(builder.source_dir)
builder.install_global_dependencies()
if %s: builder.execute_composer()
builder._execute_build_hook()
EOF
mkdir /tmp/cache
chown -R web:web /tmp/cache
echo '%s' | base64 -d | /usr/bin/python2.7 /tmp/build.py
touch /config/built
`

// appDeployCmd is the deploy command for applications.
const appDeployCmd = `
cat >/tmp/deploy.py <<EOF
#!/usr/bin/python2.7
from platformsh_gevent import patch ; patch()
import platformsh.agent
with platformsh.agent.log_and_load() as service:
	service._post_deploy()
EOF
chmod +x /tmp/deploy.py
/tmp/deploy.py
`

// serviceContainerCmd is the service container start command.
const serviceContainerCmd = `
umount /etc/hosts
umount /etc/resolv.conf
mkdir -p /run/shared /run/rpc_pipefs/nfs
cat >/tmp/fake-rpc.py <<EOF
from gevent.monkey import patch_all;
patch_all();
from gevent_jsonrpc import RpcServer;
import json;
def rootFactory(c, a):
	with open("/tmp/.ready", "w") as f: f.write("true")
	c.send(json.dumps({"jsonrpc":"2.0","result":True,"id": json.loads(c.recv(1024))["id"]}))
RpcServer(
	"/run/shared/agent.sock",
	"foo",
	root=None,
	root_factory=rootFactory
)._accepter_greenlet.get();
EOF
python /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
sleep 1
runsvdir -P /etc/service &> /tmp/runsvdir.log &
sleep 1
until [ -f /run/config.json ]; do sleep 1; done
/etc/platform/boot
exec init
`

// serviceStartCmd is the command to start a service.
const serviceStartCmd = `
until [ -f /run/config.json ]; do sleep 1; done
/etc/platform/start &
until [ -f /run/config.json ]; do sleep 1; done
`

// serviceOpenCmd is the command to open a service.
const serviceOpenCmd = `
until [ -f /tmp/.ready ]; do sleep 1; done
echo '%s' | base64 -d | /etc/platform/commands/open
`
