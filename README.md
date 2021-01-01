Platform.CC (Platform.ContextualCode)
=====================================
**By Nathan Ogden / Contextual Code**

Requirements / Installation
---------------------------

### Requirements

- Docker (See https://docs.docker.com/desktop/)

### Installation

Run the following in a Bash shell...
```curl -s https://platform-cc-releases.s3.amazonaws.com/install.sh | bash /dev/stdin```

After that use the `pcc` command to run Platform.CC.


Quick Start
-----------

### Platform.sh Cloner / Syncer

A shell script is included that allows syncing a Platform.sh environment to a local Platform.CC environment. To use it you need the SSH URL to your Platform.sh environment and the current user should have permission to access the environment over SSH.

```pcc_psh_sync <ssh_url>```

This script does not clone the repository. Make sure you do that prior and run the script inside the root directory of the project.

This script will sync the following...

- Mount directories (Rsync).
- Variables.
- MySQL databases.

### Local Project

This assumes that you have a project ready to go with all the appropiate configuration files (.platform.app.yaml, etc). Before running any `pcc` commands you should make sure that you are in the project's root directory.

1) (Optional) Set project environment variables.

    Often times environment variables are used to expose secure credientials to a project.
    You can use the `var:set` command to set environment variables.

    In order to specifically set an environment variable the variable name must be prefixed with `env:`.

    Example...
    ```
    $ pcc var:set env:ACCESS_TOKEN secret_token_here
    ```

    A common use case might be to set the 'COMPOSER_AUTH' environment variable so that Composer can run
    without user interaction being required.

    Composer auth example...
    ```
    $ pcc var:set 'env:COMPOSER_AUTH' " `cat ~/.composer/auth.json | tr -d '\n\r '` "
    ```

2) Start in the root directory of the project.

    ```
    $ pcc project:start
    ```

    This will pull all the needed Docker images and run the build commands for your application(s).

3) (Optional) Install database dumps.

    All databases will work but Platform.CC currently only has commands to interact with MySQL/MariaDB. For all other database types you will need to user `docker exec` to access them.

    For MySQL databases you can use the `pcc mysql:sql` command to execute queries.

    STDIN example...

    ```
    $ pcc mysql:sql -d main < dump.sql
    ```

    Platform.sh CLI example...

    ```
    $ platform db:dump -emaster -fdb_dump.sql
    $ pcc mysql:sql -d main < db_dump.sql
    ```

4) Run deploy hooks.

    Deploy hooks are not ran automatically. You can run your project's deploy hooks with the following command...

    ```
    $ pcc project:deploy
    ```

    This runs the deploy hooks defined in .platform.app.yaml.


Supported Languages
-------------------

Platform.CC should support all languages that are supported by Platform.sh. Please see https://docs.platform.sh/ for a list of languages.


Supported Services
------------------

Platform.CC should support all services that are supported by Platform.sh. Please see https://docs.platform.sh/ for a list of services.


Flags
-----

Some functionality is disabled by default. You can re-enable the functionality by setting flags.


**enable_cron**

`pcc project:flag:set enable_cron`
Enable cron jobs for the current project.


**enable_service_routes**

`pcc project:flag:set enable_service_routes`
Enable routes to services such as Varnish.


**enable_workers**

`pcc project:flag:set enable_workers`
Enables workers.


**enable_php_opcache**

`pcc project:flag:set enable_php_opcache`
Enables PHP Opcache.

**enable_mount_volume**

`pcc project:flag:set enable_mount_volume`
Enables mount volumes.

**enable_osx_nfs_mounts**
`pcc project:flag:set enable_osx_nfs_mounts`
Enables NFS mounts on OSX.


Options
-------

Options are like flags except that they are specific values and not just on or off.

**domain_suffix**

`pcc project:options:set domain_suffix <value>`
Set the internal route domain, default is "pcc.localtest.me." Every route gets an internal route... example... www.example.com becomes www-example-com.pcc.localtest.me.


Global Configuration
--------------------

You can create global configuration file that allows setting configuration that applies to all projects. Platform.CC looks for the global configuration file in the following locations...

```
~/.config/platform_cc.yaml
~/platform_cc.yaml
```

**Variables**

You can configure project variables that get applied to all project.

Example...
```
variables:
    env:
        COMPOSER_AUTH: '{"github-oauth":{"github.com":"SECRET_KEY_HERE"}'
```

**Router**

You can configure the ports used by the router.

Example...
```
router:
    http: 80
    https: 443
```


Platform.CC Specific Configurations
-----------------------------------

If you find that you need some configurations that are specific only to your Platform.CC projects, you can put those in a file called `.platform.app.pcc.yaml`. This should be in the same format as your `.platform.app.yaml` file.


Share Logs
----------

A script is included that will allow you to share your Platform.CC logs with ease.

`pcc_send_log`

After running the command it will give you a URL that you can share.