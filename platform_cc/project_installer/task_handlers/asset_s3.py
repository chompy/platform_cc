import os
import boto3
import datetime
import collections
import re
import difflib
import io
from .base import BaseTaskHandler

class AssetS3TaskHandler(BaseTaskHandler):

    """
    Task handler for downloading assets from S3 bucket
    to local file system.
    """

    @classmethod
    def getType(cls):
        return "asset_s3"

    def run(self):
        self.checkParams(["bucket", "from", "to"])

        # get aws creds
        awsAccessKey = self.project.variables.get("env:AWS_ACCESS_KEY_ID")
        awsSecretKey = self.project.variables.get("env:AWS_SECRET_ACCESS_KEY")
        awsRegion = self.project.variables.get("env:AWS_DEFAULT_REGION", "us-east-1")
        # ensure creds are provided
        if not awsAccessKey or not awsSecretKey:
            raise ValueError(
                "Task handler '%s' requires variables 'env:AWS_ACCESS_KEY_ID' and 'env:AWS_SECRET_ACCESS_KEY' to be set." % (
                    self.getType()
                )
            )

        # get download from
        downloadFrom = self.params.get("from")
        downloadFrom = downloadFrom.replace(
            "{PROJECT_DIRNAME}",
            os.path.basename(self.project.path)
        )
        now = datetime.datetime.now()
        downloadFrom = now.strftime(downloadFrom)

        # parse download to path
        app, downloadTo = self.parseAppPath(self.params.get("to"))

        # log that we are searching for match on s3
        self.logger.info(
            "Locate S3 asset that matches 's3://%s/%s.'" % (
                self.params.get("bucket"),
                downloadFrom,
            )
        )

        # init boto s3
        s3 = boto3.resource(
            "s3",
            aws_access_key_id=awsAccessKey,
            aws_secret_access_key=awsSecretKey,
            region_name=awsRegion
        )
        # get bucket resource
        bucket = s3.Bucket(self.params.get("bucket"))

        # get bucket objects
        bucketObjects = bucket.objects.all()
        matchKey = None
        for obj in bucketObjects:
            # perform regex match on 'downloadFrom'
            m = re.search(downloadFrom, obj.key)
            if m:
                matchKey = obj.key
                break
        if not matchKey:
            self.logger.warn("No assets found that matched expression '%s.'" % downloadFrom)
            return

        self.logger.info(
            "Found 's3://%s/%s', download to '%s:%s.'" % (
                self.params.get("bucket"),
                matchKey,
                app.getName(),
                downloadTo
            )
        )

        # download asset
        downloadTemp = io.BytesIO()
        bucket.download_fileobj(
            matchKey,
            downloadTemp
        )
        downloadTemp.seek(0)

        # move asset to app container
        app.uploadFile(downloadTemp, "/tmp/dump")
        downloadTemp.close()
        app.runCommand(
            "cd /app && mv /tmp/dump %s" % downloadTo
        )