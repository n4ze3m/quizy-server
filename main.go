package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/n4ze3m/quizy-server/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client
var authHandler *handlers.AuthHandler
var quizHandler *handlers.QuizHandler

func init() {
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	authHandler = handlers.NewAuthHandler(client.Database("quizy").Collection("users"), ctx)
	quizHandler = handlers.NewQuizHandler(client.Database("quizy").Collection("quiz"), ctx)

	log.Println("Connected react-queryto MongoDB!")
}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	// No auth middleware for public routes
	v1.POST("/login", authHandler.SignInHandler)
	v1.POST("/register", authHandler.SignUpHandler)

	// Auth middleware for user routes
	user := v1.Group("/user")
	user.Use(authHandler.AuthMiddleware())
	{
		user.GET("/list", quizHandler.GetAllUserQuizHandler)
		user.POST("/create", quizHandler.CreateQuizHandler)
	}

	router.Run()
}
