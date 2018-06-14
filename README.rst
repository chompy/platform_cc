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


Extra Options
-------------

Platform.CC contains a few additional features that are disabled by default. They can be enabled and
disabled by using the 'project:option_set' command. A project restart is required to fully enable and
disable the features. For a list of features and their current status use the 'project:options' command.

**USE_MOUNT_VOLUMES**
When enabled mount points defined in .platform.app.yaml are mounted to a Docker volume.

**ENABLE_CRON**
Enables Cron tasks as defined in .platform.app.yaml.


Missing Features
----------------

See TODO for list of features the still need to be implementd.

Currently Unsupported Functionality:

- Worker container.

Currently Unplanned Functionality:

- Non PHP applications
- Limiting app size and disk space
- Web upstream,socket_family (PHP doesn't really need this?)

