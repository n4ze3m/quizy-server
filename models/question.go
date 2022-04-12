package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Slug      string             `json:"slug" bson:"slug"`
	Title     string             `json:"title" bson:"title"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Questions []Question         `json:"questions" bson:"questions"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Question struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Question string             `json:"question" bson:"question"`
	Options  []Options          `json:"options" bson:"options"`
}

type Options struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Value    string             `json:"value" bson:"value"`
	IsAnswer bool               `json:"is_answer" bson:"is_answer"`
}

func (q *Quiz) AllFalse() {
	for _, question := range q.Questions {
		for _, option := range question.Options {
			option.IsAnswer = false
		}
	}
}

func (q *Quiz) Validate() (bool, string) {

	if strings.Trim(q.Title, " ") == "" {
		return false, "Title is required"
	}

	if len(q.Questions) == 0 {
		return false, "Quiz must have at least one question"
	}

	for _, question := range q.Questions {

		if strings.Trim(question.Question, " ") == "" {
			return false, "Question is required"
		}

		if len(question.Options) == 0 {
			return false, "Question must have at least two option"
		} else if len(question.Options) == 1 {
			return false, "Question must have at least two option"
		} else if len(question.Options) > 4 {
			return false, "Question must have at most four option"
		}

		for _, option := range question.Options {
			if strings.Trim(option.Value, " ") == "" {
				return false, "Option is required"
			}
		}
	}

	return true, ""
}

func (q *Quiz) SetID() {
	q.ID = primitive.NewObjectID()
	q.Slug = uuid.New().String()

	for i, question := range q.Questions {
		question.ID = primitive.NewObjectID()
		for j, option := range question.Options {
			option.ID = primitive.NewObjectID()
			q.Questions[i].Options[j] = option
		}
		q.Questions[i] = question
	}

}
