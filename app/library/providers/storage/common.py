import hashlib
import requests
# import uuid
import os
import tempfile


def get_file(logger, storage_root, doc_uuid, image_url):
    # local_path = "%s/%s/images/tmp/%s.%s" % (storage_root, doc_uuid, uuid.uuid4(), image_url["extension"])
    try:
        r = requests.get(image_url["url"])
        temp = tempfile.NamedTemporaryFile(delete=False, suffix=".%s" % (image_url["extension"]))
        temp.write(r.content)
        logger.info("%s => %s" % (image_url["url"], temp.name))
    except Exception as e:
        logger.exception(e)
        return '', True

    return temp.name, hash_file(temp.name), False


def create_path(storage_root, doc_uuid):
    directory = "%s/%s/images/tmp" % (storage_root, doc_uuid)
    if not os.path.exists(directory):
        os.makedirs(directory, exist_ok=True)


def hash_file(filename):
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
