package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DB NotesDatabase

type NotesDatabase interface {
	List(ctx context.Context, filter, sort string, tags []string) ([]Note, error)
	Get(ctx context.Context, id string) (Note, error)
	Create(ctx context.Context, data Note) (Note, error)
	Update(ctx context.Context, data Note) error
	Delete(ctx context.Context, id string) (Note, error)
	Tags(ctx context.Context) ([]string, error)
}

type Response struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Note struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Content  string             `bson:"content" json:"content"`
	Title    string             `bson:"title" json:"title"`
	Created  int                `bson:"created" json:"created,omitempty"`
	Updated  int                `bson:"updated" json:"updated,omitempty"`
	Tags     []string           `bson:"tags" json:"tags,omitempty"`
	Response *Response          `bson:"api_response" json:"api_response,omitempty"`
	Error    string             `bson:"error" json:"error,omitempty"`
}

type NoteNoID struct {
	Content string   `bson:"content" json:"content"`
	Title   string   `bson:"title" json:"title"`
	Created int      `bson:"created" json:"created,omitempty"`
	Updated int      `bson:"updated" json:"updated,omitempty"`
	Tags    []string `bson:"tags" json:"tags,omitempty"`
}
