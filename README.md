Platform.CC (Platform.ContextualCode)
=====================================
**By Nathan Ogden / Contextual Code**

Requirements / Installation
---------------------------

### Requirements

- Docker
    - Mac OS:
        1.  Install Version 17.09.1-ce-mac42 https://download.docker.com/mac/stable/21090/Docker.dmg. If you are unable to install that version - try to upgrade your OS.
        2.  Open Docker -> Preferences -> File Sharing and remove all directories except `/tmp`.
        3.  Open Docker -> Preferences -> Advanced and increase memory to 4.0 GB.
        4.  Install `d4m-nfs`:

                $ cd ~
                $ git clone git@github.com:IFSight/d4m-nfs.git
                $ cd d4m-nfs
                $ echo /Users:/Users > etc/d4m-nfs-mounts.txt
                $ sudo rm /etc/exports && sudo touch /etc/exports
                $ ./d4m-nfs.sh

    - Linux
        - Debian

            Follow the instructions here: https://docs.docker.com/install/linux/docker-ce/debian/

        - Ubuntu

            Follow the instructions here: https://docs.docker.com/install/linux/docker-ce/ubuntu/

        - CentOS

            Follow the instructions here: https://docs.docker.com/install/linux/docker-ce/centos/

### Installation

Run the following in a Bash shell...
`curl -s https://gitlab.com/contextualcode/platform_cc/-/raw/v2.0.x/install.sh | bash /dev/stdin`

After that use the `pcc` command to run Platform_cc.


Quick Start
-----------

### Platform.sh Cloner

TODO...provide bash script to perform this function.


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

    All databases will work but Platform_cc currently only has commands to interact with MySQL/MariaDB. For all other database types you will need to user `docker exec` to access them.

    For MySQL databases you can use the `pcc mysql:sql` command to execute queries.

    STDIN example...

    ```
    $ pcc mysql:sql -d main < dump.sql
    ```

    Platform.sh CLI example...

    ```
    $ platform db:dump -emaster -fdb_dump.sql
    $ pcc mysql:sql -d main < db_dump.sql

4) Run deploy hooks.

    Deploy hooks are not ran automatically. You can run your project's deploy hooks with the following command...

    ```
    $ pcc project:deploy
    ```

    This runs the deploy hooks defined in .platform.app.yaml. If you have multiple applications you will
    need to run this command for each application.


Supported Languages
-------------------

Platform_cc should support all languages that are supported by Platform.sh. Please see https://docs.platform.sh/ for a list of languages.


Supported Services
------------------

Platform_cc should support all services that are supported by Platform.sh. Please see https://docs.platform.sh/ for a list of services.


Options and Flags
-----------------

Some functionality is disabled by default. You can reenable the functionality by setting flags...


**enable_cron**

`pcc project:flag:set enable_cron`
Enable cron jobs for the current project.


**enable_service_routes**

`pcc project:flag:set enable_service_routes`
Enable routes to services such as Varnish.


You can also set options...

**domain_suffix**

`pcc project:options:set domain_suffix <value>`
Set the internal route domain, default is "pcc.localtest.me." Every route gets an internal route... example... www.example.com becomes www-example-com.pcc.localtest.me.


Platform.CC Specific Configurations
-----------------------------------

If you find that you need some configurations that are specific only to your Platform_cc projects, you can put those in a file called `.platform.app.pcc.yaml`. This should be in the same format as your `.platform.app.yaml` file.

This is still WIP with this version of Platform_cc. It partially works...it needs to be fixed to merge values from .platform.app.yaml better instead of doing a full override.