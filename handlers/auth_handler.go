package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthHandler(collection *mongo.Collection, ctx context.Context) *QuizHandler {
	return &AuthHandler{collection: collection, ctx: ctx}
}

func (h *AuthHandler) SignUp(c *gin.Context) {

}
