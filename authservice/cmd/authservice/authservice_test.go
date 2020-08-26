package authservice

import (
	"context"
	"testing"

	"github.com/benweissmann/memongo"
	"github.com/golang-friends/slack-clone/authservice/models"
	pb "github.com/golang-friends/slack-clone/authservice/protos/authservice"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Helper method to insert a temp user
func insertTempUser(t *testing.T, username, email, password, testdburl string) {
	t.Helper()
	// Connect and create temp user, email, and password
	models.ConnectToTestDB(testdburl)
	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	models.Db.Collection("user").InsertOne(context.Background(), models.UserInMongoDb{ID: primitive.NewObjectID(), Email: email, Username: username, Password: string(pw)})
}

func Test_authServer_Login(t *testing.T) {
	// InMemroy MongoDB
	mongoServer, err := memongo.Start("4.0.5")
	assert.NoError(t, err)
	defer mongoServer.Stop()

	// Insert our temp user
	insertTempUser(t, "incidrthreat", "incidrthreat@gmail.com", "incidrthreatpass", mongoServer.URIWithRandomDB())

	server := AuthServer{}
	_, err = server.Login(context.Background(), &pb.LoginRequest{
		User: &pb.User{
			Username: "incidrthreat",
			Email:    "incidrthreat@gmail.com",
			Password: "incidrthreatpass",
		}})
	if err != nil {
		t.Error("1. An error was returned: ", err.Error())
	}

	_, err = server.Login(context.Background(), &pb.LoginRequest{
		User: &pb.User{
			Username: "incidrthreat@gmail.com",
			Password: "notincidrthreatpass",
		}})
	if err == nil {
		t.Error("2. Password should be incorrect")
	}

	_, err = server.Login(context.Background(), &pb.LoginRequest{
		User: &pb.User{
			Username: "incidrthreat",
			Password: "incidrthreatpass",
		}})
	if err != nil {
		t.Error("3. An error was returned: ", err.Error())
	}

	// Wrong Email with correct password
	_, err = server.Login(context.Background(), &pb.LoginRequest{
		User: &pb.User{
			Email:    "notincidrthreat@gmail.com",
			Password: "incidrthreatpass",
		},
	},
	)
	if err == nil {
		t.Error("4. Email should be incorrect and return an error.")
	}
}

func Test_authServer_UsernameUsed(t *testing.T) {
	// InMemroy MongoDB
	mongoServer, err := memongo.Start("4.0.5")
	assert.NoError(t, err)
	defer mongoServer.Stop()

	// Insert our temp user
	insertTempUser(t, "incidrthreat", "incidrthreat@gmail.com", "incidrthreatpass", mongoServer.URIWithRandomDB())

	server := AuthServer{}

	res, err := server.UsernameUsed(context.Background(), &pb.UsernameUsedRequest{Username: "incidrthreat"})
	// 1. Server responded with an error
	if err != nil {
		t.Error("1. An error was returned: ", err.Error())
	}
	// 2. Our username exists, If the server responded [false] which it shouldn't.. the test fails
	if !res.GetUsed() {
		t.Error("2. Username is used, should have returned true")
	}

	res, err = server.UsernameUsed(context.Background(), &pb.UsernameUsedRequest{Username: "testUser"})
	// 3. Server responded with an error
	if err != nil {
		t.Error("3. An error was returned: ", err.Error())
	}
	// 3. Our username does not exist, If the server responded [true] which it shouldn't.. the test fails
	if res.GetUsed() {
		t.Error("3. User name is not used, should have returned false")
	}

}

// TODO

// Test EmailUsed
// Test Register
