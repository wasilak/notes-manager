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

func NewMongoDB(ctx context.Context) (*MongoDB, error) {
	ctx, span := common.TracerCmd.Start(ctx, "NewMongoDB")
	mongoDB := &MongoDB{}
	err := mongoDB.setup(ctx)
	if err != nil {
		return nil, err
	}
	span.End()
	return mongoDB, nil
}

func (d *MongoDB) setup(ctx context.Context) error {
	ctx, span := common.TracerCmd.Start(ctx, "setup")
	mongoConnectionString := os.Getenv("MONGO_CONNECTION_STRING")
	if mongoConnectionString == "" {
		mongoConnectionString = "mongodb://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASS") + "@" + os.Getenv("MONGO_HOST")
	}

	ctxCancel, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ctx, spanClient := common.TracerCmd.Start(ctx, "client")
	client, err := mongo.Connect(ctxCancel, options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		span.End()
		return err
	}
	spanClient.End()

	d.client = client

	ctx, spanDatabase := common.TracerCmd.Start(ctx, "database")
	d.db = client.Database("notes")
	spanDatabase.End()

	d.setupIndices(ctx)
	span.End()
	return nil
}

func (d *MongoDB) setupIndices(ctx context.Context) {
	ctx, span := common.TracerCmd.Start(ctx, "setupIndices")
	curIndices, err := d.db.Collection("notes").Indexes().List(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	textIndexName := "content_text_title_text"
	textIndexCreated := false

	for curIndices.Next(ctx) {
		var index bson.M
		err := curIndices.Decode(&index)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
		}

		if index["name"] == textIndexName {
			textIndexCreated = true
			break
		}
	}

	if !textIndexCreated {
		slog.InfoContext(ctx, "Creating text index")
		opts := options.Index().SetWeights(bson.M{"title": 10, "content": 1})
		_, err := d.db.Collection("notes").Indexes().CreateOne(ctx,
			mongo.IndexModel{
				Keys:    bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
				Options: opts,
			})
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
		}
	}
	span.End()
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

func (d *MongoDB) List(ctx context.Context, filter, sort string, tags []string) ([]Note, error) {
	ctx, span := common.TracerWeb.Start(ctx, "List")
	searchParams := bson.M{}
	sortParams := bson.D{}

	var opts *options.FindOptions

	tags = deleteEmpty(tags)

	if len(tags) > 0 {
		searchParams["tags"] = bson.M{"$all": tags}
	}

	filterElements := deleteEmpty(strings.Split(filter, " "))
	if len(filterElements) > 0 {
		orQuery := make([]bson.M, 0)
		for _, v := range filterElements {
			orQuery = append(orQuery, bson.M{"title": bson.M{"$regex": v, "$options": "i"}})
			orQuery = append(orQuery, bson.M{"content": bson.M{"$regex": v, "$options": "i"}})
		}
		searchParams["$or"] = orQuery
	}

	sortField := "updated"
	sortOrder := -1

	if len(sort) > 0 {
		sortTmp := strings.Split(sort, ":")
		sortField = sortTmp[0]
		sortOrder = 1
		if sortTmp[1] == "desc" {
			sortOrder = -1
		}
	}

	sortParams = append(sortParams, bson.E{Key: sortField, Value: sortOrder})

	opts = options.Find().SetSort(sortParams)

	cur, err := d.db.Collection("notes").Find(ctx, searchParams, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	docs := []Note{}
	for cur.Next(ctx) {
		var doc Note
		if err := cur.Decode(&doc); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, err
		}
		docs = append(docs, doc)
	}

	span.End()

	return docs, nil
}

func (d *MongoDB) Get(ctx context.Context, id string) (Note, error) {
	ctx, span := common.TracerWeb.Start(ctx, "Get")
	var note Note

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Note{}, err
	}

	err = d.db.Collection("notes").FindOne(ctx, bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		return Note{}, err
	}

	span.End()
	return note, nil
}

func (d *MongoDB) Create(ctx context.Context, data Note) (Note, error) {
	ctx, span := common.TracerWeb.Start(ctx, "Create")
	data.ID = primitive.NewObjectID()

	result, err := d.db.Collection("notes").InsertOne(ctx, data)
	if err != nil {
		slog.ErrorContext(ctx, "error", "result", result)
		return Note{}, err
	}

	span.End()
	return data, nil
}

func (d *MongoDB) Update(ctx context.Context, note Note) error {
	ctx, span := common.TracerWeb.Start(ctx, "Update")
	filter := bson.M{"_id": note.ID}

	replacement := NoteNoID{
		Content: note.Content,
		Title:   note.Title,
		Created: note.Created,
		Updated: note.Updated,
		Tags:    note.Tags,
	}

	result, err := d.db.Collection("notes").ReplaceOne(ctx, filter, replacement)
	if err != nil {
		slog.ErrorContext(ctx, "error", "result", result)
		slog.ErrorContext(ctx, err.Error())
	}

	span.End()
	return err
}

func (d *MongoDB) Delete(ctx context.Context, id string) (Note, error) {
	ctx, span := common.TracerWeb.Start(ctx, "Delete")
	note, err := d.Get(ctx, id)
	if err != nil {
		return Note{}, err
	}

	filter := bson.M{"_id": note.ID}

	result, err := d.db.Collection("notes").DeleteOne(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "error", "result", result)
		return Note{}, err
	}

	span.End()
	return note, nil
}

func (d *MongoDB) Tags(ctx context.Context) ([]string, error) {
	ctx, span := common.TracerWeb.Start(ctx, "Tags")

	ctx, spanTags := common.TracerWeb.Start(ctx, "getTags")
	cur, err := d.db.Collection("notes").Distinct(ctx, "tags", bson.M{})
	if err != nil {
		return nil, err
	}
	spanTags.End()

	_, spanTagsConvert := common.TracerWeb.Start(ctx, "convertTags")
	tags := []string{}

	for _, tag := range cur {
		if tag != nil && len(tag.(string)) > 0 {
			tags = append(tags, tag.(string))
		}
	}
	spanTagsConvert.End()

	span.End()
	return tags, nil
}
