package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (u *User) HashPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(bytes)
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}


type AuthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
