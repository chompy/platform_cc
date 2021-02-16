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
# INIT
usermod -u %d app
groupmod -g %d app
usermod -u %d web
groupmod -g %d web
umount /etc/hosts
umount /etc/resolv.conf
mkdir -p /run/shared /run/rpc_pipefs/nfs
# MOUNT TMP
mount -o user_xattr --bind /mnt/data/_tmp /tmp
# CLEAN UP TMP
rm -rf /tmp/cache
rm -rf /tmp/sessions
rm -rf /tmp/log
rm -rf /tmp/nginx
rm -f /tmp/*.py
rm -f /tmp/*.log
rm -f /tmp/.ready*
# FAKE RPC
cat >/tmp/fake-rpc.py <<EOF
from gevent.monkey import patch_all;
patch_all();
from gevent_jsonrpc import RpcServer;
import json;
def rootFactory(c, a):
	with open("/tmp/.ready2", "w") as f: f.write("true")
	c.send(json.dumps({"jsonrpc":"2.0","result":True,"id": json.loads(c.recv(1024))["id"]}))
RpcServer(
	"/run/shared/agent.sock",
	"foo",
	root=None,
	root_factory=rootFactory
)._accepter_greenlet.get();
EOF
# INIT SERVICE
until [ -f /run/config.json ]; do sleep 1; done
rm -rf /etc/service/*
python /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
runsvdir -P /etc/service &> /tmp/runsvdir.log &
# PERMISSIONS
chown -R web:web /run
chown -R web:web /mnt/data/_tmp
chown -R web:web /tmp
chmod -R 0777 /tmp
chmod -R 0777 /mnt/data/_tmp
mkdir -p /run/sshd
chown -R root:root /run/sshd
chmod -R -rwx /run/sshd
rm -f /run/rsa_hostkey
# BOOT
/etc/platform/boot
sleep 5
touch /tmp/.ready1
exec init
`

// appBuildCmd is the build command for applications.
const appBuildCmd = `
until [ -f /run/config.json ]; do sleep 1; done
until [ -f /tmp/.ready1 ]; do sleep 1; done
until [ -f /tmp/.ready2 ]; do sleep 1; done
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
mkdir -p /tmp/cache
chown -R web:web /tmp/cache
chmod -R 0777 /tmp/cache
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
	with open("/tmp/.ready2", "w") as f: f.write("true")
	c.send(json.dumps({"jsonrpc":"2.0","result":True,"id": json.loads(c.recv(1024))["id"]}))
RpcServer(
	"/run/shared/agent.sock",
	"foo",
	root=None,
	root_factory=rootFactory
)._accepter_greenlet.get();
EOF
until [ -f /run/config.json ]; do sleep 1; done
rm -rf /etc/service/*
python /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
runsvdir -P /etc/service &> /tmp/runsvdir.log &
mkdir -p /run/sshd
chown -R root:root /run/sshd
chmod -R -rwx /run/sshd
rm -f /run/rsa_hostkey
/etc/platform/boot
sleep 5
touch /tmp/.ready1
exec init
`

// serviceStartCmd is the command to start a service.
const serviceStartCmd = `
until [ -f /run/config.json ]; do sleep 1; done
until [ -f /tmp/.ready1 ]; do sleep 1; done
until [ -f /tmp/.ready2 ]; do sleep 1; done
/etc/platform/start &
`

// serviceOpenCmd is the command to open a service.
const serviceOpenCmd = `
until [ -f /run/config.json ]; do sleep 1; done
until [ -f /tmp/.ready1 ]; do sleep 1; done
until [ -f /tmp/.ready2 ]; do sleep 1; done
echo '%s' | base64 -d | /etc/platform/commands/open
`
