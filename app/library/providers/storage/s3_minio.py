import os
import logging
from minio import Minio
from .common import get_file, create_path
from datetime import timedelta


class Storage():

    bucket_name = os.getenv("S3_BUCKET", "notes")
    app_root = os.getcwd()
    storage_root = "%s/storage" % (app_root)

    def __init__(self):
        if not os.path.exists(self.storage_root):
            os.makedirs(self.storage_root, exist_ok=True)

        self.logger = logging.getLogger("uvicorn.error")

        self.client = Minio(
            os.getenv("MINIO_ADDRESS", ""),
            access_key=os.getenv("MINIO_ACCESS_KEY", ""),
            secret_key=os.getenv("MINIO_SECRET_KEY", ""),
            region=os.getenv("MINIO_REGION_NAME", ""),
        )

        # Make bucket if not exist.
        found = self.client.bucket_exists(self.bucket_name)
        if not found:
            self.client.make_bucket(self.bucket_name)
        else:
            self.logger.info("Bucket {} already exists".format(self.bucket_name))

    def get_files(self, doc_uuid, image_urls):
        for item in image_urls:
            create_path(self.storage_root, doc_uuid)
            local_path, file_hash, error = get_file(self.logger, self.storage_root, doc_uuid, item["original"])

            if not error:
                filename = "%s/storage/images/%s.%s" % (doc_uuid, file_hash, item["original"]["extension"])

                try:
                    self.client.fput_object(self.bucket_name, filename, local_path)

                    item["replacement"] = "/storage/%s" % (filename)
                except Exception as e:
                    self.logger.exception(e)

    def cleanup(self, doc_uuid):
        # Remove a prefix recursively.
        delete_object_list = map(lambda x: DeleteObject(x.object_name), self.client.list_objects(self.bucket_name, "%s/" % (doc_uuid), recursive=True))
        errors = self.client.remove_objects("my-bucket", delete_object_list)
        for error in errors:
            self.logger.error("error occured when deleting object", error)

    def get_object(self, filename, expiration=20):
        # response = self.client.fget_object(self.bucket_name, filename, filename)
        response = self.client.presigned_get_object(self.bucket_name, filename, expires=timedelta(hours=expiration))
        self.logger.info(filename)
        self.logger.info(response)

        return response
