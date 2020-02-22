import os
import logging
import boto3
from botocore.exceptions import ClientError
from .common import get_file, create_path


class Storage():

    bucket_name = os.getenv("S3_BUCKET", "notes")
    app_root = os.getcwd()
    storage_root = "%s/storage" % (app_root)

    def __init__(self):
        if not os.path.exists(self.storage_root):
            os.makedirs(self.storage_root, exist_ok=True)

        self.s3_client = boto3.client(
            's3',
            aws_access_key_id=os.getenv("S3_ID", ""),
            aws_secret_access_key=os.getenv("S3_SECRET", ""),
            region_name=os.getenv("S3_REGION", "")
        )

        self.s3_resource = boto3.resource(
            's3',
            aws_access_key_id=os.getenv("S3_ID", ""),
            aws_secret_access_key=os.getenv("S3_SECRET", ""),
            region_name=os.getenv("S3_REGION", "")
        )

        self.bucket = self.s3_resource.Bucket(self.bucket_name)

        self.logger = logging.getLogger("api")

    def get_files(self, doc_uuid, image_urls):
        for item in image_urls:
            create_path(self.storage_root, doc_uuid)
            local_path, file_hash, error = get_file(self.logger, self.storage_root, doc_uuid, item["original"])

            if not error:
                filename = "%s/storage/images/%s.%s" % (doc_uuid, file_hash, item["original"]["extension"])

                try:
                    self.s3_client.upload_file(local_path, self.bucket_name, filename)
                    item["replacement"] = "/storage/%s" % (filename)
                except ClientError as e:
                    self.logger.exception(e)

    def cleanup(self, doc_uuid):
        self.bucket.objects.filter(Prefix="%s/" % (doc_uuid)).delete()

    def get_object(self, filename, expiration=20):
        response = self.s3_client.generate_presigned_url(
            'get_object',
            Params={
                'Bucket': self.bucket_name,
                'Key': filename
            },
            ExpiresIn=expiration
        )

        return response
