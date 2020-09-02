package models

import (
	"fmt"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang-friends/slack-clone/authservice/protos/authservice"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGetToken(t *testing.T) {
	id := primitive.NewObjectID()

	userInMongoDb := UserInMongoDb{
		ID:       id,
		Username: "username",
		Password: "password",
		Email:    "email",
		Admin:    false,
	}

	tokenString := userInMongoDb.GetToken([]byte("jwt-secret"), time.Minute)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("jwt-secret"), nil
	})

	assert.NoError(t, err)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user authservice.User
		err = proto.Unmarshal([]byte(claims["user"].(string)), &user)
		assert.NoError(t, err)

		assert.Empty(t,
			cmp.Diff(&authservice.User{
				Id:       id.Hex(),
				Username: "username",
				// Password is not set
				Email: "email",
				Admin: false,
			}, &user,
				protocmp.Transform(),
			))
	} else {
		assert.Fail(t, "expected tokenClaims is ok")
	}
}
