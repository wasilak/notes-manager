import os


class Storage():

    # def setup(self):
    #     self.tinydb = TinyDB('./db.json')
    #     self.table = self.tinydb.table('notes', cache_size=0)
    #     self.Note = Query()

    def create_path(self, directory):
        if not os.path.exists(directory):
            os.makedirs(directory, exist_ok=True)
