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
until [ -f /tmp/.ready1 ]; do sleep 1; done
/usr/bin/python2.7 /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
runsvdir -P /etc/service &> /tmp/runsvdir.log &
exec init
`

// appInitCmd is the initalization command for applications.
const appInitCmd = `
# INIT
usermod -u %d app
groupmod -g %d app
usermod -u %d web
groupmod -g %d web
umount /etc/hosts
umount /etc/resolv.conf
mkdir -p /run/shared /run/rpc_pipefs/nfs
# MOUNT TMP
mkdir -p /mnt/data/.tmp
mount -o user_xattr --bind /mnt/data/.tmp /tmp
# MOUNT TMP/CACHE
mkdir -p /var/pcc_global/cache
mkdir -p /tmp/cache
mount -o user_xattr --bind /var/pcc_global/cache/ /tmp/cache
# CLEAN UP TMP
rm -rf /tmp/sessions
rm -rf /tmp/log
rm -rf /tmp/nginx
rm -f /tmp/*.py
rm -f /tmp/*.log
rm -f /tmp/.ready*
# CLEAN UP RUN
rm -rf /run/ssh/*
rm -rf /run/sshd/*
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
# CLEAN UP SERVICE
rm -rf /etc/service/*
# PERMISSIONS
chown -R web:web /run
chown -R web /tmp
chmod -R 0755 /tmp
mkdir -p /run/sshd
chown -R root:root /run/sshd
chmod -R -rwx /run/sshd
chown -Rf root:root /run/ssh/id
chmod -Rf -rwx /run/ssh/id
rm -f /run/rsa_hostkey
# BOOT
/etc/platform/boot
chown -R web /tmp
chmod -R 0755 /tmp
touch /tmp/.ready1
`

// appBuildCmd is the build command for applications.
const appBuildCmd = `
until [ -f /tmp/.ready1 ]; do sleep 1; done
timeout 1m bash -c 'until [ -f /tmp/.ready2 ]; do sleep 1; done'
touch /tmp/.ready2
chown -R web /tmp
chmod -R 0755 /tmp
# UPDATE COMPOSER
if [ -f /usr/bin/composer ]; then
	composer self-update -q -n
fi
# NOTE: we don't want the builder method move_source_directory to execute in PCC
# TODO this could break in the future....
if [ -f /etc/platform/flavor.d/composer.py ]; then
	sed -i '20,25d' /etc/platform/flavor.d/composer.py
fi
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
builder._build()
if builder.execute_build_hook:
	builder._execute_build_hook()
	builder.prepare_mounts()
EOF
chown -R web /tmp/cache
chmod -R 0755 /tmp/cache
chown -R web /app
echo '%s' | base64 -d | /usr/bin/python2.7 /tmp/build.py
chown -R web /tmp
chmod -R 0755 /tmp
chown -R web /app
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
until [ -f /tmp/.ready1 ]; do sleep 1; done
/usr/bin/python2.7 /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
runsvdir -P /etc/service &> /tmp/runsvdir.log &
exec init
`

// serviceInitCmd is the command to initalize a service.
const serviceInitCmd = `
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
rm -rf /etc/service/*
mkdir -p /run/sshd
chown -R app:app /run
chown -R root:root /run/sshd
chmod -R -rwx /run/sshd
chown -Rf root:root /run/ssh/id
chmod -Rf -rwx /run/ssh/id
rm -f /run/rsa_hostkey
/etc/platform/boot
touch /tmp/.ready1
`

// serviceStartCmd is the command to start a service.
const serviceStartCmd = `
until [ -f /tmp/.ready1 ]; do sleep 1; done
timeout 1m bash -c 'until [ -f /tmp/.ready2 ]; do sleep 1; done'
touch /tmp/.ready2
/etc/platform/start &
`

// serviceOpenCmd is the command to open a service.
const serviceOpenCmd = `
until [ -f /tmp/.ready1 ]; do sleep 1; done
timeout 1m bash -c 'until [ -f /tmp/.ready2 ]; do sleep 1; done'
touch /tmp/.ready2
echo '%s' | base64 -d | /etc/platform/commands/open
`
