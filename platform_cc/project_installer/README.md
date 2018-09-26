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
          bucket: database-dumps
          from: "project.sql.gz"
          to: "app:/app/project.sql.gz"

        - type: mysql_import
          from: app:/app/project.sql.gz
          to: mysqldb:main
          delete_dump: true

        - type: command
          to: app
          command: |
                php app/console cache:clear


## Available Tasks

**asset_s3**

Fetch asset from Amazon S3 bucket.

Parameters...

- bucket :: Bucket to fetch from.
- from :: Path to asset in bucket to fetch.
- to :: Application and path inside container to dump asset to. (<app_name>:<path_to>)


**mysql_import**

Import a MySQL database dump.

Parameters...

- from :: Application and path inside container to fetch dump from. (<app_name>:<path_to>)
- to :: Service and database name to import dump to. (<service_name>:<database_name>)
- delete_dump :: If set to true then the original dump file will be deleted upon completion.