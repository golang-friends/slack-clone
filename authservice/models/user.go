package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-friends/slack-clone/authservice/configs"

	"github.com/dgrijalva/jwt-go"
)

// NilUser ...
var NilUser UserInMongoDb

// UserInMongoDb ...
type UserInMongoDb struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Email    string             `bson:"email"`
	Admin    bool               `bson:"admin"`
}

// GetToken ...
func (u UserInMongoDb) GetToken() string {
	var conf configs.Configuration
	byteSlc, _ := json.Marshal(u)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string(byteSlc),
	})

	tokenString, _ := token.SignedString([]byte(conf.JWTSecret))
	return tokenString
}
