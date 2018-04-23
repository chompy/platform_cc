Platform.CC (Platform.ContextualCode)
=====================================
**By Nathan Ogden / Contextual Code**


Tool for provisioning apps with Docker based on Platform.sh's .platform.app.yaml spec.


Quick Start
-----------

This assumes that you have a project ready to go with all the appropiate configuration files (.platform.app.yaml, etc).

1) Install Composer credientials.

    You can use the "var:set" command to set environment variables. This can be use to
    setup Composer credientials for your project. The below command is an example
    of how one might copy their .composer/auth.json in to their project.

        $ platform_cc var:set env:COMPOSER_AUTH `cat ~/.composer/auth.json | tr '\r' ' ' |  tr '\n' ' ' | sed 's/ \{3,\}/ /g' | sed 's/   / /g'`

2) Start in the root directory of the project.

        $ platform_cc project:start

    This will pull all the needed Docker images and run the build commands for your application(s).
    
3) Setup Databases and Other Services

    You should now install databases and setup other services.

    You can use the 'mysql:sql' command to run SQL queries and gain access to the MySQL console.

4) Run deploy hooks.

        $ platform_cc application:deploy   

    This runs the deploy hooks defined in .platform.app.yaml. If you have multiple applications you will
    need to run this command for each application.

        $ platform_cc application:deploy --name=application_name


More On Variables
-----------------

The 'platform_cc var:set' command allows you to set variables that are exposed to your
project. All the variables are exposed to your application via the PROJECT_VARIABLES
environment variable. Additionally any variable prefixed with "env:" will be set as an
environment variable.


Mount Volumes
-------------

By default all application mount points are ignored. You can configure Platform CC to create a Docker volume to
bind to all mount points by setting the project config value 'use_mount_volumes' in .pcc_project.json.

In the future I'd like to add the ability to fully mount the entire application to a volume. Additionally it'd
be nice to mount the application code as read only to more closely simulate a Platform.sh environment.


Missing Features
----------------

See TODO for list of features the still need to be implementd.

Currently Unsupported Functionality:

- Cron tasks

Currently Unplanned Functionality:

- Non PHP applications
- Workers
- Limiting app size and disk space
- Web upstream,socket_family (PHP doesn't really need this?)

