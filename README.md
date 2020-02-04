Platform.CC (Platform.ContextualCode)
=====================================
**By Nathan Ogden / Contextual Code**

Tool for provisioning apps with Docker based on Platform.sh's .platform.app.yaml spec.


Requirements / Installation
---------------------------

### Requirements

- Python 3.6+ (2.7+ has worked in the past but is no longer supported.)
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

Make sure you have Python 3 and Pip installed then run one of the following set of commands...

    $ python3 -m pip install git+https://gitlab.com/contextualcode/platform_cc    

OR

    $ cd ~
    $ git clone https://gitlab.com/contextualcode/platform_cc.git
    $ cd platform_cc
    $ pip install -r requirements.txt
    $ sudo python setup.py install

Quick Start
-----------

### Platform.sh Cloner

The following assumes you already have a project up and running on Platform.sh and that you have added your SSH key.

You can clone a project straight from Platform.sh. Follow these steps...

1) Visit https://accounts.platform.sh/user/api-tokens to create an API token.
2) Run the following command...

    ```
    $ platform_cc platform_sh:login <API_TOKEN>
    ```
3) Add your SSH private key to Platform.CC by running the following command...

    ```
    $ platform_cc platform_sh:set_ssh -p ~/.ssh/id_rsa
    ```
4) Obtain the Platform.sh project ID. You can obtain a list of your project ids from the Platform.sh CLI tool or visit the project dashboard and copy it from the URL... `https://console.platform.sh/<USERNAME>/<PROJECT_ID>`

5) Run the following command to start the cloning process...

    ```
    $ platform_cc platform_sh:clone <PROJECT_ID>
    ```

### Local Project

This assumes that you have a project ready to go with all the appropiate configuration files (.platform.app.yaml, etc). Before running any `platform_cc` commands you should make sure that you are in the project's root directory.

1) (Optional) Set project environment variables.

    Often times environment variables are used to expose secure credientials to a project.
    You can use the `var:set` command to set environment variables.

    In order to specifically set an environment variable the variable name must be prefixed with `env:`.

    Example...
    ```
    $ platform_cc var:set env:ACCESS_TOKEN secret_token_here
    ```

    A common use case might be to set the 'COMPOSER_AUTH' environment variable so that Composer can run
    without user interaction being required.

    Composer auth example...
    ```
    $ platform_cc var:set -g 'env:COMPOSER_AUTH' " `cat ~/.composer/auth.json | tr -d '\n\r '` "
    ```

    Also note the `-g` option. It allows you to set global variables that affect all projects ran under the same user.


2) Start in the root directory of the project.

    ```
    $ platform_cc project:start
    ```

    This will pull all the needed Docker images and run the build commands for your application(s).
    
3) (Optional) Install database dumps.

    As of writting this MySQL is the only supported database. You can use the `mysql:sql` command to access the MySQL interactive shell and run queries. You can also pass queries in through STDIN.

    STDIN example...

    ```
    $ cat db_dump.sql | platform_cc mysql:sql -d main
    ```

    Platform.sh CLI example...

    ```
    $ platform db:dump -emaster -fdb_dump.sql
    $ cat db_dump.sql | platform_cc mysql:sql -d main
    ```

4) Run deploy hooks.

    Deploy hooks are not ran automatically. You can run your project's deploy hooks with the following command...

    ```
    $ platform_cc application:deploy
    ```

    This runs the deploy hooks defined in .platform.app.yaml. If you have multiple applications you will
    need to run this command for each application.

    ```
    $ platform_cc application:deploy --name application_name
    ```


Supported Languages
-------------------

Platform.CC was primarily designed to aid with PHP development, however additional language support has been and can be added. Here is a list of the currently supported languages and their versions...

- PHP 5.6
- PHP 7.0
- PHP 7.1
- PHP 7.2
- PHP 7.3
- PHP 7.4
- Go 1.11
- Python 3.7


Supported Services
------------------

- MySQL 5.5
- MySQL 10.0
- MySQL 10.1
- MySQL 10.2
- Memcached 1.4
- Redis 2.8
- Redis 3.0
- Redis 3.2
- RabbitMQ 3.5

Non Platform.sh supported services...

- Minio (Object store, like Amazon S3)
- Athenapdf (HTML to PDF api)
- Docker (Allows creation of custom services via Docker images)


Upgrading To New Platform.CC Version
------------------------------------

It is likely that new base application images will be available when a new Platform.CC version
is released. There are cases where old application images will no longer be compatible with the
new version of Platform.CC. In that event you can upgrade your application by running the following
commands...

```
$ platform_cc application:pull
$ platform_cc application:build
```

If your project contains multiple applications you will likely need to run the command for each of them
using the --name argument to specify the application.


Extra Options
-------------

Platform.CC contains a few additional features that are disabled by default. They can be enabled and
disabled by using the 'project:option_set' command. A project restart is required to fully enable and
disable the features. For a list of features and their current status use the 'project:options' command.

- **USE_MOUNT_VOLUMES**

    When enabled mount points defined in .platform.app.yaml are mounted to a Docker volume. Setting this to `true` is important for performance on Macs.

- **ENABLE_CRON**
    
    Enables Cron tasks as defined in .platform.app.yaml.

- **ENABLE_SERVICE_ROUTES**

    Enables routing through services. This is mostly used to enable Varnish. Since Platform.CC is a development platform first this option is disabled by default.


Platform.CC Specific Configurations
-----------------------------------

If you find that you need some configurations that are specific only to your Platform.CC projects, you can put those in a file called `.platform.app.pcc.yaml`. This should be in the same format as your `.platform.app.yaml` file.

A simple example scenario is to enforce 'dev' mode for your project when ran through Platform.CC. If you added the following to `.platform.app.pcc.yaml`...

```
variables:
    env:
        SYMFONY_ENV: dev
```

Then the environment variables `SYMFONY_DEV` would always get set to `dev` in Platform.CC but not on Platform.sh. This also works for `services.yaml` and `routes.yaml` which would be named `services.pcc.yaml` and `routes.pcc.yaml` respectively.


Missing Features
----------------

Platform.sh adds new features all the time. Keeping up with them was never the goal of Platform.CC...rather the goal is to focus on functionality we at Contextual Code need for our day to day development.

### Known Missing Features

- Workers
