from tinydb import TinyDB, Query


class Db:

    def __init__(self):
        self.tinydb = TinyDB('./db.json')
        self.table = self.tinydb.table('notes', cache_size=0)
        self.Note = Query()

    def list(self, filter):
        if len(filter) > 0:
            return self.table.search((self.Note.title.matches(filter)) | (self.Note.content.matches(filter)))
        return self.table.all()

    def get(self, id):
        return self.table.get(self.Note.id == id)

    def create(self, id, data):
        self.update(id, data)

    def update(self, id, data):
        self.table.upsert(data, self.Note.id == id)

    def delete(self, id):
        item = self.get(id)
        self.table.remove(self.Note.id == id)
        return item
