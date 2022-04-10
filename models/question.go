package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Questions []Question         `json:"questions" bson:"questions"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Question struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Question  string             `json:"question" bson:"question"`
	Options   []Options          `json:"options" bson:"options"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Options struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Value     string             `json:"value" bson:"value"`
	IsAnswer  bool               `json:"is_answer" bson:"is_answer"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
