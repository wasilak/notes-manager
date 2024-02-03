package db

import (
	"context"
	"os"
	"strings"

	"log/slog"

	"github.com/wasilak/notes-manager/libs/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDB(ctx context.Context) (*MongoDB, error) {
	ctx, span := common.TracerCmd.Start(ctx, "MongoDB_NewMongoDB")
	defer span.End()

	mongoDB := &MongoDB{}
	err := mongoDB.setup(ctx)
	if err != nil {
		common.HandleError(ctx, err)
		return nil, err
	}
	return mongoDB, nil
}

func (d *MongoDB) setup(ctx context.Context) error {
	ctx, span := common.TracerCmd.Start(ctx, "MongoDB_setup")
	defer span.End()

	mongoConnectionString := os.Getenv("MONGO_CONNECTION_STRING")
	if mongoConnectionString == "" {
		mongoConnectionString = "mongodb://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASS") + "@" + os.Getenv("MONGO_HOST")
	}

	ctx, spanClient := common.TracerCmd.Start(ctx, "MongoDB_client")

	// connect to MongoDB
	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(mongoConnectionString)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		common.HandleError(ctx, err)
		return err
	}
	spanClient.End()

	d.client = client

	ctx, spanDatabase := common.TracerCmd.Start(ctx, "MongoDB_database")
	d.db = client.Database("notes")
	spanDatabase.End()

	err = d.setupIndices(ctx)
	if err != nil {
		common.HandleError(ctx, err)
		return err
	}

	return nil
}

func (d *MongoDB) setupIndices(ctx context.Context) error {
	ctx, span := common.TracerCmd.Start(ctx, "MongoDB_setupIndices")
	defer span.End()

	curIndices, err := d.db.Collection("notes").Indexes().List(ctx)
	if err != nil {
		common.HandleError(ctx, err)
		return err
	}

	textIndexName := "content_text_title_text"
	textIndexCreated := false

	for curIndices.Next(ctx) {
		var index bson.M
		err := curIndices.Decode(&index)
		if err != nil {
			common.HandleError(ctx, err)
			return err
		}

		if index["name"] == textIndexName {
			textIndexCreated = true
			break
		}
	}

	if !textIndexCreated {
		slog.DebugContext(ctx, "Creating text index")
		opts := options.Index().SetWeights(bson.M{"title": 10, "content": 1})
		_, err := d.db.Collection("notes").Indexes().CreateOne(ctx,
			mongo.IndexModel{
				Keys:    bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
				Options: opts,
			})
		if err != nil {
			common.HandleError(ctx, err)
			return err
		}
	}
	return nil
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
	ctx, span := common.TracerWeb.Start(ctx, "MongoDB_List")
	defer span.End()

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
		common.HandleError(ctx, err)
		return nil, err
	}
	defer cur.Close(ctx)

	docs := []Note{}
	for cur.Next(ctx) {
		var doc Note
		if err := cur.Decode(&doc); err != nil {
			common.HandleError(ctx, err)
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (d *MongoDB) Get(ctx context.Context, id string) (Note, error) {
	ctx, span := common.TracerWeb.Start(ctx, "MongoDB_Get")
	defer span.End()

	var note Note

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		common.HandleError(ctx, err)
		return Note{}, err
	}

	err = d.db.Collection("notes").FindOne(ctx, bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		common.HandleError(ctx, err)
		return Note{}, err
	}

	return note, nil
}

func (d *MongoDB) Create(ctx context.Context, data Note) (Note, error) {
	ctx, span := common.TracerWeb.Start(ctx, "MongoDB_Create")
	defer span.End()

	data.ID = primitive.NewObjectID()

	_, err := d.db.Collection("notes").InsertOne(ctx, data)
	if err != nil {
		common.HandleError(ctx, err)
		return Note{}, err
	}

	return data, nil
}

func (d *MongoDB) Update(ctx context.Context, note Note) error {
	ctx, span := common.TracerWeb.Start(ctx, "MongoDB_Update")
	defer span.End()

	filter := bson.M{"_id": note.ID}

	replacement := NoteNoID{
		Content: note.Content,
		Title:   note.Title,
		Created: note.Created,
		Updated: note.Updated,
		Tags:    note.Tags,
	}

	_, err := d.db.Collection("notes").ReplaceOne(ctx, filter, replacement)
	if err != nil {
		common.HandleError(ctx, err)
		common.HandleError(ctx, err)
	}

	return err
}

func (d *MongoDB) Delete(ctx context.Context, id string) (Note, error) {
	ctx, span := common.TracerWeb.Start(ctx, "MongoDB_Delete")
	defer span.End()

	note, err := d.Get(ctx, id)
	if err != nil {
		common.HandleError(ctx, err)
		return Note{}, err
	}

	filter := bson.M{"_id": note.ID}

	_, err = d.db.Collection("notes").DeleteOne(ctx, filter)
	if err != nil {
		common.HandleError(ctx, err)
		return Note{}, err
	}

	return note, nil
}

func (d *MongoDB) Tags(ctx context.Context) ([]string, error) {
	ctx, span := common.TracerWeb.Start(ctx, "MongoDB_Tags")
	defer span.End()

	ctx, spanTags := common.TracerWeb.Start(ctx, "MongoDB_getTags")
	cur, err := d.db.Collection("notes").Distinct(ctx, "tags", bson.M{})
	if err != nil {
		common.HandleError(ctx, err)
		return nil, err
	}
	spanTags.End()

	_, spanTagsConvert := common.TracerWeb.Start(ctx, "MongoDB_convertTags")
	tags := []string{}

	for _, tag := range cur {
		if tag != nil && len(tag.(string)) > 0 {
			tags = append(tags, tag.(string))
		}
	}
	spanTagsConvert.End()

	return tags, nil
}
