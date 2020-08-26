package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang-friends/slack-clone/authservice/protos/authservice"
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
func (u UserInMongoDb) GetToken(jwtSecret []byte, expiry time.Duration) string {
	byteSlc, _ := proto.Marshal(u.ToUserProto())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": string(byteSlc),
		"exp":  time.Now().Add(expiry).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

// ToUserProto is a convenient method to convert to UserProto.
func (u UserInMongoDb) ToUserProto() *authservice.User {
	return &authservice.User{
		Id:       u.ID.Hex(),
		Username: u.Username,
		Email:    u.Email,
		Admin:    u.Admin,
	}
}
