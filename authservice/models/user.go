package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
)

// NilUser ...0
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
func (u UserInMongoDb) GetToken(jwtSecret []byte, expiry int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":  u.Username,
		"admin": u.Admin,
		"exp":   time.Now().Add(time.Minute * time.Duration(expiry)).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}
