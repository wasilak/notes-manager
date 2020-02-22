import os
import shutil
import logging
from .common import get_file, create_path


class Storage():

    app_root = os.getcwd()
    storage_root = "%s/storage" % (app_root)

    def __init__(self):
        if not os.path.exists(self.storage_root):
            os.makedirs(self.storage_root, exist_ok=True)

        self.logger = logging.getLogger("api")

    def create_path(self, doc_uuid):
        directory = "%s/%s/images/tmp" % (self.storage_root, doc_uuid)
        if not os.path.exists(directory):
            os.makedirs(directory, exist_ok=True)

    def get_files(self, doc_uuid, image_urls):
        for item in image_urls:
            create_path(self.storage_root, doc_uuid)
            local_path, file_hash, error = get_file(self.logger, self.storage_root, doc_uuid, item["original"])

            if not error:
                replacement_path = "%s/%s/images/%s.%s" % (self.storage_root, doc_uuid, file_hash, item["original"]["extension"])
                os.rename(local_path, replacement_path)
                item["replacement"] = replacement_path.replace(self.app_root, "")

    def cleanup(self, doc_uuid):
        self.remove_dir(doc_uuid)

    def remove_dir(self, doc_uuid):
        directory = "%s/%s" % (self.storage_root, doc_uuid)
        if os.path.exists(directory):
            shutil.rmtree(directory)
