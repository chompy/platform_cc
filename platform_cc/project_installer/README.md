# Platform.CC - Project Installer

The project installer, 'platform_cc project:install,' is a subcommand in Platform.CC. It's purpose it to execute a series of task from an 'install.pcc.yaml' file that should be present in the root of the project you intend to install. The purpose of this installer command is to provide private credientials and install assets (such as database dumps) to the project. This 'install.pcc.yaml' file is not to replace the 'platform.app.yaml' file.


## install.pcc.yaml Syntax

**vars**

This section allows you to set project variables with key/value pairs. This is the equivalent of manually running 'platform_cc var:set.' You can set environment variables by prefixing the key with 'env:.'

Example...

    vars:
        "env:AWS_ACCESS_KEY_ID": "AWS_ACCESS_KEY_HERE"
        "env:AWS_SECRET_ACCESS_KEY": "AWS_SECRET_HERE"


**tasks**

Example...

    tasks:
        - type: asset_s3
          from: "database-dumps/project.sql.gz"
          to: "app:/app/project.sql.gz"

        - type: mysql_import
          from: app:/app/project.sql.gz
          to: mysqldb:main
          delete_dump: true
          condition: ["app", "[ -f /app/project.sql.gz ]"]

        - type: command
          to: app
          command: |
                php app/console cache:clear


## Task Run Condition

Every task can execute a shell command that dictates if the task should be ran or not. If the shell command
returns a non zero exit code then the task will be skipped.

Any one of the following syntaxes will work...

    condition: ["<app_name>", "<shell_command>"]

    condition: "<shell_command>"

    condition:
        application: "<app_name>"
        command: "<shell_command>"

If <app_name> is omitted then the first available application will be used.


## Available Tasks

**asset_s3**

Fetch asset from Amazon S3 bucket.

Parameters...

- from :: Bucket and path to asset to fetch. (<bucket_name>/<path_to>)
- to :: Application and path inside container to dump asset to. (<app_name>:<path_to>)


**mysql_import**

Import a MySQL database dump.

Parameters...

- from :: Application and path inside container to fetch dump from. (<app_name>:<path_to>)
- to :: Service and database name to import dump to. (<service_name>:<database_name>)
- delete_dump :: If set to true then the original dump file will be deleted upon completion.


**shell**

Run a shell command inside an application container.

Parameters...

- to :: Name of application, uses first found application if not provided.
- command :: Shell command to run.


**rsync**

Rsync remote asset(s) to application directory.

Parameters...

- from :: Remote path to asset(s). (user@server.com:/path/to)
- to :: Application and path inside container to copy asset(s) to. (<app_name>:<path_to>)