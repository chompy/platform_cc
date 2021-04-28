Platform.CC (Platform.ContextualCode)
=====================================
**By Nathan Ogden / Contextual Code**

Requirements / Installation
---------------------------

### Requirements

- Docker (See https://docs.docker.com/desktop/)

### Installation

Run the following in a Bash shell...
```
curl -L -s https://platform.cc/releases/install.sh | bash /dev/stdin
```

After that use the `pcc` command to run Platform.CC.

### Development Release

You can also install the latest development build with...
```
curl -L -s https://platform.cc/releases/install.sh | bash /dev/stdin -d
```


Quick Start
-----------

### Platform.sh Sync

If you wish to copy an environment from Platform.sh you can use the `platformsh:sync` command. To use this command you must first obtain a Platform.sh API token. You can do this by using the `platformsh:login` command.

Steps...

1) Run `pcc platformsh:login` and follow the instructions.
2) From the project root directory run `pcc platformsh:sync`. This will take some time.

This command will sync the following...

- Mount directories (Rsync). (Skipped with `--skip-mounts`.)
- Variables. (Skipped with `--skip-variables`.)
- MySQL/Postgres databases. (Skipped with `--skip-databases`.)

On first sync Platform.CC will generate a new SSH key, it can take some time for Platform.sh to accept the new key.

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
    $ pcc var:set 'env:COMPOSER_AUTH' "$(cat ~/.composer/auth.json | tr -d '\n\r ')"
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

5) Open in web browser.

    Visit http://localhost in your browser for a list of routes.


Supported Languages
-------------------

Platform.CC should support all languages that are supported by Platform.sh. Please see https://docs.platform.sh/ for a list of languages.


Supported Services
------------------

Platform.CC should support all services that are supported by Platform.sh. Please see https://docs.platform.sh/ for a list of services.


Valid SSL Certificates
----------------------

The main router generates a certficate CA that you can have your browser trust to get valid SSL certificates. To get the CA certificate you can run...

```
pcc router:certificates:dump
```

Flags, Options, and Variables
-----------------------------

You can set various flags, options, and project variables through the command line. Doing so will set the value for the current project. If you use the `--global` (`-g`) option then the value will be set globally for all projects.

```
pcc project:flag:set <NAME>                 # set flag
pcc project:flag:set --unset <NAME>         # unset flag
pcc project:flag:set --off <NAME>           # turn flag off (override global)
pcc project:flag:set -g <NAME>              # set global flag
pcc project:flag:set -g --unset <NAME>      # unset global flag
pcc project:flag:list                       # list all flags

pcc project:option:set <NAME> <VALUE>       # set option
pcc project:option:set -g <NAME> <VALUE>    # set global option
pcc project:option:reset <NAME>             # unset option
pcc project:option:list                     # list all options

pcc variable:set <NAME> <VALUE>             # set a project variable
pcc variable:set -g <NAME> <VALUE>          # set a global project variable
pcc variable:get <NAME>                     # get a project variable
pcc variable:delete <NAME>                  # delete a project variable
pcc variable:delete -g <NAME>               # delete a global project variable
pcc variable:list                           # list all project variables
```

## Flags

Flags enables or disable functionality. Some functionality is disabled by default (such as workers and cron jobs) and can be re-enabled with flags.

### **Lists of available flags...**

### enable_cron
Enable cron jobs for the current project.

### enable_service_routes
Enable routes to services such as Varnish.

### enable_workers
Enables workers.

### enable_php_opcache
Enables PHP Opcache.

### enable_osx_nfs_mounts
Enables NFS mounts on OSX.

### disable_yaml_overrides
Disables Platform.CC specific YAML overrides (.platform.app.pcc.yaml and services.pcc.yaml).

### disable_auto_commit
Disables the auto commit of application containers when a project is started.


## Options

Options are like flags except that they are specific values and not just on or off.

### domain_suffix
Set the internal route domain, default is "pcc.localtest.me." Every route gets an internal route... example... www.example.com becomes www-example-com.pcc.localtest.me.

### mount_strategy [none|symlink|volume]
Set the strategy for how Platform.CC will deal with mounts.
- None (default)...do nothing, no mounts will be established. This will suffice in a lot of cases if you have no need for recursive mounts.
- Symlink...create a mount directory in the root of the project and use symlinks to map the mounts to the destinations. This should create a mount structure similar to Platform.sh and will support recursive mounts, it could be destructive to your project directory if you previously had data in the mount directories.
- Volume...use a container volume as the mount directory and bind the destination directories to the container volume. This should function closest to how it would on Platform.sh but files in the mounted directories will not be accessible outside of the application container.


## SSH

Platform.cc generates its own SSH keypair on first load. This is used inside application containers and to perform syncs with Platform.sh.

You can find the Platform.cc generated key in... 

### Private
`~/.config/platformcc/pcc_ssh_private`

### Public
`~/.config/platformcc/pcc_ssh_public`


Platform.CC Specific Configurations
-----------------------------------

If you find that you need some configurations that are specific only to your Platform.CC projects, you can put those in a file called `.platform.app.pcc.yaml`. This should be in the same format as your `.platform.app.yaml` file.


Slots
-----

When you start a project you can specify a 'save slot' to use with the 'slot' or 's' option. This lets you start the project with different storage locations for the services in your project. An example use case is loading a different MySQL database without losing the original. The slot value must be an integer.

Examples...
```
pcc project:start -s 1
pcc project:start --slot 2
```

### Delete Slot
```
pcc project:slot:delete 2
```

### Copy Slot
```
pcc project:slot:copy 1 2
```


Container Commit
----------------

By default Platform.CC makes a commit of the application containers once the build hooks are completed. This saves time as the build hooks won't have to be ran on every start up. However, it does appear to cause some issues in a few cases. We have included ways to bypass this functionality both with the 'disable_auto_commit' flag and with the `--no-commit` option on the `project:start` command.

You can force a re-build and re-commit when you run `project:start` or `project:restart` with the `--rebuild` flag.


Self Update
-----------

You can update Platform.cc to the latest version with the `update` command.


Share Logs
----------

The `sharelog` command will upload the log for the current project and provide a URL to share it. This is useful for reporting bugs to a Platform.cc developer.

