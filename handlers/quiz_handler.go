package handlers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n4ze3m/quizy-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuizHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewQuizHandler(collection *mongo.Collection, ctx context.Context) *QuizHandler {
	return &QuizHandler{collection: collection, ctx: ctx}
}

func (q *QuizHandler) CreateQuizHandler(c *gin.Context) {
	id := c.Keys["user"].(string)
	var quiz models.Quiz

	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}


	inValid, msg := quiz.Validate()

	if !inValid {
		c.JSON(400, gin.H{"error": msg})
		return
	}

	quiz.SetID()
	quiz.UserID = id
	quiz.CreatedAt = time.Now()

	if _, err := q.collection.InsertOne(q.ctx, quiz); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Quiz created successfully", "slug": quiz.Slug})
}


func (q *QuizHandler) GetAllUserQuizHandler(c *gin.Context) {
	id := c.Keys["user"].(string)

	var quizzes []models.Quiz = []models.Quiz{}

	cursor, err := q.collection.Find(q.ctx, bson.M{"user_id": id})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(q.ctx, &quizzes); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"quizzes": quizzes})
}


func (q *QuizHandler) GetQuizBySlugPublic(c *gin.Context) {
	slug := c.Param("slug")

	var quiz models.Quiz

	cursor, err := q.collection.Find(q.ctx, bson.M{"slug": slug})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !cursor.Next(q.ctx) {
		c.JSON(404, gin.H{"error": "Quiz not found"})
		return
	}

	if err := cursor.Decode(&quiz); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	quiz.AllFalse()

	c.JSON(200, gin.H{"quiz": quiz})
}