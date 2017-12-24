Platform.CC (Platform.ContextualCode)
=====================================
**By Nathan Ogden / Contextual Code**


Tool for provisioning apps with Docker based on Platform.sh's .platform.app.yaml spec.


Quick Start
-----------

This assumes that you have a project ready to go with all the appropiate configuration files (.platform.app.yaml, etc).

1) Start in the root directory of the project.

        $ platform_cc project:start

    This will pull all the needed Docker images so it can take a while.

2) Install SSH key and Composer auth files. (Needed for Composer projects.)

    The 'var:set' command will allow you to install credientials needed by Composer for your project. Below
    are examples that will allow you to install credientials from your local machine.

        $ platform_cc var:set env:COMPOSER_AUTH `cat ~/.composer/auth.json | tr '\r' ' ' |  tr '\n' ' ' | sed 's/ \{3,\}/ /g' | sed 's/   / /g'`
        $ platform_cc var:set project:ssh_key `cat ~/.ssh/id_rsa | base64 -w 0`
        $ platform_cc var:set project:known_hosts `cat ~/.ssh/known_hosts | base64 -w 0`

    Note that 'project:ssh_key' and 'project:known_hosts' are base64 encoded.

3) Provision project.
    
        $ platform_cc project:provision

    This setups the Docker containers for your project. It will install dependencies for your project as well as run build hooks (composer install, etc).

4) Setup services. (Deploy MySql databases, etc).

    You should now setup your services such as your database.

    You can use the 'mysql:sql' command to run SQL queries and gain access to the MySQL console.

4) Deploy project.

        $ platform_cc project:deploy   

    This runs the deploy hooks defined in .platform.app.yaml.


More On Variables
-----------------

The 'platform_cc var:set' command allows you to set variables that are exposed to your project. There a several builtin variables that are used to configure your environment.

===================== =========================================================================== ==============
Variable Name         Description                                                                 Base64 Encoded
--------------------- --------------------------------------------------------------------------- --------------
project:ssh_key       SSH key (.ssh/id_rsa). Composer needs to checkout private repos.            Yes
project:known_hosts   SSH known_hosts file (.ssh/known_hosts). Needed for non interactive.        Yes
project:domains       Comma delimited list of domain names, replaces {DEFAULT} in routes.yml.     No
project:hosts_file    Additional entries in hosts file. JSON encoded, key=hostname, value=ip.     Yes
project:web_uid       UID to use for user 'web' when provisioning project.                        No
env:*                 The prefix "env:" sets an environment variable inside the app container(s). No


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

