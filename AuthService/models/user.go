package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-friends/slack-clone/AuthService/configs"

	"github.com/dgrijalva/jwt-go"
)

// NilUser ...
var NilUser User

// User ...
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Email    string             `bson:"email"`
	Admin    bool               `bson:"admin"`
}

// GetToken ...
func (u User) GetToken() string {
	var conf configs.Configuration
	byteSlc, _ := json.Marshal(u)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string(byteSlc),
	})

	tokenString, _ := token.SignedString([]byte(conf.JWTSecret))
	return tokenString
}
