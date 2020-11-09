package api

// appContainerCmd - application container start command
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
RpcServer(
	"/run/shared/agent.sock",
	"foo",
	root=None,
	root_factory=lambda c,a: c.send(json.dumps({"jsonrpc":"2.0","result":True,"id": json.loads(c.recv(1024))["id"]})))._accepter_greenlet.get();
EOF
python /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
runsvdir -P /etc/service &> /tmp/runsvdir.log &
chown -R web:web /run
until [ -f /run/config.json ]; do sleep 1; done
/etc/platform/boot
exec init
`

// appOpenCmd - command to open application
const appOpenCmd = `
/etc/platform/start &
sleep 1
echo '%s' | base64 -d | /etc/platform/commands/open
`

// appBuildScript - build script for applications
const appBuildScript = `#!/usr/bin/python2.7 -u
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
builder.execute_composer()
builder._execute_build_hook()
`

// appBuildCmd - build command for applications
const appBuildCmd = `
chmod +x /opt/build.py
mkdir /tmp/cache
chown -R web:web /tmp/cache
echo '%s' | base64 -d | /opt/build.py
`
