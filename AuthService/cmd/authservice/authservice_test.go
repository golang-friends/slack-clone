package authservice

import (
	"context"
	"testing"

	"github.com/golang-friends/slack-clone/AuthService/models"
	pb "github.com/golang-friends/slack-clone/AuthService/protos"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Test_authServer_Login(t *testing.T) {
	// Connect and create temp user, email, and password
	models.ConnectToTestDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte("incidrthreatpass"), bcrypt.DefaultCost)
	models.Db.Collection("user").InsertOne(context.Background(), models.User{ID: primitive.NewObjectID(), Email: "incidrthreat@gmail.com", Username: "incidrthreat", Password: string(pw)})

	server := AuthServer{}
	_, err := server.Login(context.Background(), &pb.LoginRequest{Username: "incidrthreat", Email: "incidrthreat@gmail.com", Password: "incidrthreatpass"})
	if err != nil {
		t.Error("1. An error was returned: ", err.Error())
	}

	_, err = server.Login(context.Background(), &pb.LoginRequest{Username: "incidrthreat@gmail.com", Password: "notincidrthreatpass"})
	if err == nil {
		t.Error("2. Error was nil")
	}

	_, err = server.Login(context.Background(), &pb.LoginRequest{Username: "incidrthreat", Password: "incidrthreatpass"})
	if err != nil {
		t.Error("3. An error was returned: ", err.Error())
	}
}

// TODO

// Test UsernameUsed
// Test EmailUsed
// Test Register
