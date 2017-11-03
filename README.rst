Platform.CC (Platform.ContextualCode)
=====================================
**By Nathan Ogden / Contextual Code**


Tool for provisioning apps with Docker based on Platform.sh's .platform.app.yaml spec.


Quick Start
-----------

This assumes that you have a project ready to go with all the appropiate configuration files (.platform.app.yaml, etc).

1) Start in the root directory of the project.

        $ platform_cc.py project:start

    This will provision all the Docker containers needed for your project. It can take a while on the first run.

2) Install SSH key and Composer auth files. (Needed for Composer projects.)

    The 'var:set' command will allow you to install credientials needed by Composer for your project. Below
    are examples that will allow you to install credientials from your local machine.

        $ platform_cc.py var:set env:COMPOSER_AUTH `cat ~/.composer/auth.json`
        $ platform.cc.py var:set project:ssh_key `cat ~/.ssh/id_rsa | base64 -w 0`
        $ platform.cc.py var:set project:known_hosts `cat ~/.ssh/known_hosts | base64 -w 0`

    Note that 'project:ssh_key' and 'project:known_hosts' are base64 encoded.

3) Build project.
    
        $ platform_cc.py project:build

    This will run 'composer install' as well as run the build hooks defined in .platform.app.yaml.

4) Setup services. (Deploy MySql databases, etc).

    You should now setup your services such as your database.

    You can use the 'mysql:sql' command to run SQL queries and gain access to the MySQL console.

4) Deploy project.

        $ platform_cc.py project:deploy   

    This runs the deploy hooks defined in .platform.app.yaml.


Missing Features
----------------

See TODO for list of features the still need to be implementd.

Currently Unsupported Functionality:

- Main router only partially implemented
- Cron tasks

Currently Unplanned Functionality:

- Non PHP applications
- Workers
- Limiting app size and disk space
- Web upstream,socket_family (PHP doesn't really need this?)

