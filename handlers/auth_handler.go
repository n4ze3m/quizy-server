package handlers

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/n4ze3m/quizy-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

type Claims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler(collection *mongo.Collection, ctx context.Context) *AuthHandler {
	return &AuthHandler{collection: collection, ctx: ctx}
}

// SignInHandler handles the login request
func (h *AuthHandler) SignInHandler(c *gin.Context) {
	var userInput models.User
	var user models.User

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if userInput.Username == "" || userInput.Password == "" {
		c.JSON(400, gin.H{"error": "Missing required fields"})
		return
	}

	u, err := h.collection.Find(h.ctx, bson.M{"username": userInput.Username})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !u.Next(h.ctx) {
		c.JSON(400, gin.H{"error": "Username does not exist"})
		return
	}

	err = u.Decode(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !user.CheckPasswordHash(userInput.Password) {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	// expiration time will be 7 days from now
	// this is not a good method
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		ID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	})

}

// SignUpHandler handles the signup request
func (h *AuthHandler) SignUpHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.Name == "" || user.Username == "" || user.Password == "" {
		c.JSON(400, gin.H{"error": "Missing required fields"})
		return
	}

	u, err := h.collection.Find(h.ctx, bson.M{"username": user.Username})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if u.Next(h.ctx) {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.HashPassword()
	_, err = h.collection.InsertOne(h.ctx, user)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

// auth middleware
func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenValue := ctx.Request.Header.Get("Authorization")
		currentUser := &models.User{}
		if tokenValue == "" {
			ctx.JSON(401, gin.H{"error": "Missing token"})
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenValue, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			id, err := primitive.ObjectIDFromHex(claims.ID)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				ctx.Abort()
				return
			}
			user, err := handler.collection.Find(handler.ctx, bson.M{"_id": id})
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				ctx.Abort()
				return
			}

			if !user.Next(handler.ctx) {
				ctx.JSON(401, gin.H{"error": "Invalid token"})
				ctx.Abort()
				return
			}
			user.Decode(currentUser)
			ctx.Set("user", currentUser.ID.Hex())
		} else {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}
	}
}
