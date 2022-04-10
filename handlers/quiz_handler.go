package handlers

import (
	"context"

	// "github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuizHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewQuizHandler(collection *mongo.Collection, ctx context.Context) *QuizHandler {
	return &QuizHandler{collection: collection, ctx: ctx}
}