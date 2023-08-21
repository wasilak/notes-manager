package db

import "go.mongodb.org/mongo-driver/bson/primitive"

var DB NotesDatabase

type NotesDatabase interface {
	List(filter, sort string, tags []string) ([]Note, error)
	Get(id string) (Note, error)
	Create(data Note) (Note, error)
	Update(data Note) error
	Delete(id string) (Note, error)
	Tags() ([]string, error)
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
	Score    int                `bson:"_score" json:"_score,omitempty"`
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
