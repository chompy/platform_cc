Platform.CC (Platform.ContextualCode)
=====================================
**By Nathan Ogden / Contextual Code**

Tool for provisioning apps with Docker based on Platform.sh's .platform.app.yaml spec.


Requirements / Installation
---------------------------

### Requirements

- Python 2.7+
- Docker
    - Mac:
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

    - Linux:

        Follow the instructions here https://github.com/docker/docker-install.
        You might have to add yourself to the docker user group:

            $ sudo addgroup --system docker && sudo adduser $USER docker && newgrp docker

### Installation

Install Python 2.7 and Pip if needed.

    $ pip install git+https://gitlab.com/contextualcode/platform_cc    

OR

    $ cd ~
    $ git clone https://gitlab.com/contextualcode/platform_cc.git
    $ cd platform_cc
    $ pip install -r requirements.txt
    $ sudo python setup.py install


Quick Start
-----------

This assumes that you have a project ready to go with all the appropiate configuration files (.platform.app.yaml, etc).

1) Install Composer credentials, ssh keys, etc.

    You can use the "var:set" command to set environment variables. This can be use to
    setup Composer credientials for your project. The below command is an example
    of how one might copy their .composer/auth.json in to their project.


        $ platform_cc var:set 'env:COMPOSER_AUTH' " `cat ~/.composer/auth.json | tr -d '\n\r '` "

        $ platform_cc var:set project:ssh_key `cat ~/.ssh/id_rsa | base64 -w 0`

        $ platform_cc var:set project:known_hosts `cat ~/.ssh/known_hosts | base64 -w 0`

    (Use `base64 -b 0` on Macs)

2) Start in the root directory of the project.

        $ platform_cc project:start

    This will pull all the needed Docker images and run the build commands for your application(s).
    
3) Setup Databases and Other Services

    You should now install databases and setup other services.

    You can use the 'mysql:sql' command to run SQL queries and gain access to the MySQL console.

    For example:

    $ platform db:dump -emaster -fdb_dump.sql

    $ cat db_dump.sql | platform_cc mysql:sql -dmain

4) Run deploy hooks.

        $ platform_cc application:deploy   

    This runs the deploy hooks defined in .platform.app.yaml. If you have multiple applications you will
    need to run this command for each application.

        $ platform_cc application:deploy --name=application_name


Supported Languages
-------------------

Platform.CC was primarily designed to aid with PHP development, however additional language support has been and can be added. Here is a list of the currently supported languages and their versions...

- PHP 5.4
- PHP 5.6
- PHP 7.0
- PHP 7.1
- PHP 7.2
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

        $ platform_cc application:pull
        $ platform_cc application:build

If your project contains multiple applications you will likely need to run the command for each of them
using the --name argument to specify the application.


More On Variables
-----------------

The 'platform_cc var:set' command allows you to set variables that are exposed to your
project. All the variables are exposed to your application via the PROJECT_VARIABLES
environment variable. Additionally any variable prefixed with "env:" will be set as an
environment variable.


Extra Options
-------------

Platform.CC contains a few additional features that are disabled by default. They can be enabled and
disabled by using the 'project:option_set' command. A project restart is required to fully enable and
disable the features. For a list of features and their current status use the 'project:options' command.

**USE_MOUNT_VOLUMES**
When enabled mount points defined in .platform.app.yaml are mounted to a Docker volume. Setting this to `true` is important for performance on Macs.

**ENABLE_CRON**
Enables Cron tasks as defined in .platform.app.yaml.


Custom config
-------------

If you find that you need some variables that are specific only to your Platform.CC projects, you can put those in a file called `.platform.app.pcc.yaml`. This should be in the same format as your `.platform.app.yaml` file.

For example, if you wanted to have the environment variable `$SYMFONY_ENV` set to `dev`, you could set it with `var:set`:

    $ platform_cc var:set 'env:SYMFONY_ENV' 'dev'

But this would have to be ran each time you restarted the project. If you wanted `$SYMFONY_ENV` to always be `dev` when in Platform.CC, you can create `.platform.app.pcc.yaml` file with contents:

    variables:
        env:
            SYMFONY_ENV: dev

In this way, you can have variables and settings that are only and automatically set in your local development environments. And importantly, it uses the same syntax as your `.platform.app.yaml` files.


Missing Features
----------------

See TODO for list of features the still need to be implementd.

Currently Unsupported Functionality:

- Worker container.

Currently Unplanned Functionality:

- Non PHP applications (Go 1.11 and Python 3.7 are currently supported however)
- Limiting app size and disk space
- Web upstream,socket_family (PHP doesn't really need this?)

