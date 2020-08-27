package authservice

import (
	"context"
	"testing"

	"github.com/golang-friends/slack-clone/AuthService/models"
	pb "github.com/golang-friends/slack-clone/AuthService/protos/authservice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Helper method to insert a temp user
func insertTempUser(t *testing.T, username, email, password string) {
	t.Helper()
	// Connect and create temp user, email, and password
	models.ConnectToTestDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	models.Db.Collection("user").InsertOne(context.Background(), models.UserInMongoDb{ID: primitive.NewObjectID(), Email: email, Username: username, Password: string(pw)})
}

func Test_authServer_Login(t *testing.T) {

	// Insert our temp user
	insertTempUser(t, "incidrthreat", "incidrthreat@gmail.com", "incidrthreatpass")

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

func Test_authServer_UsernameUsed(t *testing.T) {

	// Insert our temp user
	insertTempUser(t, "incidrthreat", "incidrthreat@gmail.com", "incidrthreatpass")

	server := AuthServer{}

	res, err := server.UsernameUsed(context.Background(), &pb.UsernameUsedRequest{Username: "incidrthreat"})
	// 1. Server responded with an error
	if err != nil {
		t.Error("1. An error was returned: ", err.Error())
	}
	// 2. Our username exists, If the server responded [true] which it shouldn't.. the test fails
	if !res.GetUsed() {
		t.Error("2. Username is used, should have returned true")
	}

	res, err = server.UsernameUsed(context.Background(), &pb.UsernameUsedRequest{Username: "testUser"})
	// 3. Server responded with an error
	if err != nil {
		t.Error("3. An error was returned: ", err.Error())
	}
	// 3. Our username does not exist, If the server responded [false] which it shouldn't.. the test fails
	if res.GetUsed() {
		t.Error("3. User name is not used, should have returned false")
	}

}

// TODO

// Test EmailUsed
// Test Register
