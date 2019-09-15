import os
import wget
import uuid
import hashlib
import shutil


class Storage():

    app_root = os.getcwd()
    storage_root = "%s/storage" % (app_root)

    def setup(self):
        if not os.path.exists(self.storage_root):
            os.makedirs(self.storage_root, exist_ok=True)

    def create_path(self, doc_uuid):
        directory = "%s/%s/images/tmp" % (self.storage_root, doc_uuid)
        if not os.path.exists(directory):
            os.makedirs(directory, exist_ok=True)

    def get_files(self, doc_uuid, image_urls):
        for item in image_urls:
            local_path = self.get_file(doc_uuid, item["original"])
            file_hash = self.hash_file(local_path)
            replacement_path = "%s/%s/images/%s.%s" % (self.storage_root, doc_uuid, file_hash, item["original"]["extension"])
            os.rename(local_path, replacement_path)
            item["replacement"] = replacement_path.replace(self.app_root, "")

    def get_file(self, doc_uuid, image_url):
        local_path = "%s/%s/images/tmp/%s.%s" % (self.storage_root, doc_uuid, uuid.uuid4(), image_url["extension"])
        wget.download(image_url["url"], local_path)

        return local_path

    def hash_file(self, filename):
        # make a hash object
        h = hashlib.sha1()

        # open file for reading in binary mode
        with open(filename, 'rb') as file:

            # loop till the end of the file
            chunk = 0
            while chunk != b'':
                # read only 1024 bytes at a time
                chunk = file.read(1024)
                h.update(chunk)

        # return the hex representation of digest
        return h.hexdigest()

    def cleanup(self, doc_uuid):
        self.remove_dir(doc_uuid)

    def remove_dir(self, doc_uuid):
        directory = "%s/%s" % (self.storage_root, doc_uuid)
        if os.path.exists(directory):
            shutil.rmtree(directory)
