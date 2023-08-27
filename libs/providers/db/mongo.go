package db

import (
	"context"
	"os"
	"strings"
	"time"

	"log/slog"

	"github.com/wasilak/notes-manager/libs/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDB() (*MongoDB, error) {
	mongoDB := &MongoDB{}
	err := mongoDB.setup()
	if err != nil {
		return nil, err
	}
	return mongoDB, nil
}

func (d *MongoDB) setup() error {
	mongoConnectionString := os.Getenv("MONGO_CONNECTION_STRING")
	if mongoConnectionString == "" {
		mongoConnectionString = "mongodb://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASS") + "@" + os.Getenv("MONGO_HOST")
	}

	ctx, cancel := context.WithTimeout(common.CTX, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		return err
	}

	d.client = client
	d.db = client.Database("notes")

	d.setupIndices()

	return nil
}

func (d *MongoDB) setupIndices() {
	curIndices, err := d.db.Collection("notes").Indexes().List(common.CTX)
	if err != nil {
		slog.ErrorContext(common.CTX, err.Error())
	}

	textIndexName := "content_text_title_text"
	textIndexCreated := false

	for curIndices.Next(common.CTX) {
		var index bson.M
		err := curIndices.Decode(&index)
		if err != nil {
			slog.ErrorContext(common.CTX, err.Error())
		}

		if index["name"] == textIndexName {
			textIndexCreated = true
			break
		}
	}

	if !textIndexCreated {
		slog.InfoContext(common.CTX, "Creating text index")
		opts := options.Index().SetWeights(bson.M{"title": 10, "content": 1})
		_, err := d.db.Collection("notes").Indexes().CreateOne(common.CTX,
			mongo.IndexModel{
				Keys:    bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
				Options: opts,
			})
		if err != nil {
			slog.ErrorContext(common.CTX, err.Error())
		}
	}
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func (d *MongoDB) List(filter, sort string, tags []string) ([]Note, error) {
	searchParams := bson.M{}
	otherParams := bson.M{}
	sortParams := bson.D{}

	tags = deleteEmpty(tags)

	if len(tags) > 0 {
		searchParams["tags"] = bson.M{"$all": tags}
	}

	if len(filter) > 0 {
		searchParams["$text"] = bson.M{"$search": filter}
		otherParams["score"] = bson.M{"$meta": "textScore"}
		sortParams = append(sortParams, bson.E{Key: "score", Value: bson.M{"$meta": "textScore"}})
	}

	if len(sort) > 0 {
		sortTmp := strings.Split(sort, ":")
		sortOrder := 1
		if sortTmp[1] == "desc" {
			sortOrder = -1
		}
		sortParams = append(sortParams, bson.E{Key: sortTmp[0], Value: sortOrder})
	}

	if len(sortParams) == 0 {
		// Workaround to make searches without sort work
		sortParams = append(sortParams, bson.E{Key: "$natural", Value: 1})
	}

	opts := options.Find().SetSort(sortParams)

	cur, err := d.db.Collection("notes").Find(common.CTX, searchParams, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(common.CTX)

	docs := []Note{}
	for cur.Next(common.CTX) {
		var doc Note
		if err := cur.Decode(&doc); err != nil {
			slog.ErrorContext(common.CTX, err.Error())
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (d *MongoDB) Get(id string) (Note, error) {
	var note Note

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Note{}, err
	}

	err = d.db.Collection("notes").FindOne(common.CTX, bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		return Note{}, err
	}
	return note, nil
}

func (d *MongoDB) Create(data Note) (Note, error) {
	data.ID = primitive.NewObjectID()

	result, err := d.db.Collection("notes").InsertOne(common.CTX, data)
	if err != nil {
		slog.ErrorContext(common.CTX, "error", result)
		return Note{}, err
	}
	return data, nil
}

func (d *MongoDB) Update(note Note) error {
	filter := bson.M{"_id": note.ID}

	replacement := NoteNoID{
		Content: note.Content,
		Title:   note.Title,
		Created: note.Created,
		Updated: note.Updated,
		Tags:    note.Tags,
	}

	result, err := d.db.Collection("notes").ReplaceOne(common.CTX, filter, replacement)
	if err != nil {
		slog.ErrorContext(common.CTX, "error", result)
		slog.ErrorContext(common.CTX, err.Error())
	}
	return err
}

func (d *MongoDB) Delete(id string) (Note, error) {
	note, err := d.Get(id)
	if err != nil {
		return Note{}, err
	}

	filter := bson.M{"_id": note.ID}

	result, err := d.db.Collection("notes").DeleteOne(common.CTX, filter)
	if err != nil {
		slog.ErrorContext(common.CTX, "error", result)
		return Note{}, err
	}
	return note, nil
}

func (d *MongoDB) Tags() ([]string, error) {
	cur, err := d.db.Collection("notes").Distinct(common.CTX, "tags", bson.M{})
	if err != nil {
		return nil, err
	}

	tags := []string{}

	for _, tag := range cur {
		if tag != nil && len(tag.(string)) > 0 {
			tags = append(tags, tag.(string))
		}
	}
	return tags, nil
}
