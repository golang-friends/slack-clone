package authservice

import (
	"context"
	"regexp"
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

	// Username, email, and password correct
	_, err = server.Login(context.Background(), &pb.LoginRequest{Username: "incidrthreat", Email: "incidrthreat@gmail.com", Password: "incidrthreatpass"})
	if err != nil {
		t.Error("1. An error was returned: ", err.Error())
	}

	// Username with wrong password
	_, err = server.Login(context.Background(), &pb.LoginRequest{Username: "incidrthreat", Password: "notincidrthreatpass"})
	if err == nil {
		t.Error("2. Password should be incorrect")
	}

	// Email login with correct password
	_, err = server.Login(context.Background(), &pb.LoginRequest{Email: "incidrthreat@gmail.com", Password: "incidrthreatpass"})
	if err != nil {
		t.Error("3. An error was returned: ", err.Error())
	}

	// Wrong Email with correct password
	_, err = server.Login(context.Background(), &pb.LoginRequest{Email: "notincidrthreat@gmail.com", Password: "incidrthreatpass"})
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

func Test_authServer_EmailUsed(t *testing.T) {
	// InMemroy MongoDB
	mongoServer, err := memongo.Start("4.0.5")
	assert.NoError(t, err)
	defer mongoServer.Stop()

	// Insert our temp user
	insertTempUser(t, "incidrthreat", "incidrthreat@gmail.com", "incidrthreatpass", mongoServer.URIWithRandomDB())
	server := AuthServer{}

	res, err := server.EmailUsed(context.Background(), &pb.EmailUsedRequest{Email: "incidrthreat@gmail.com"})
	// 1. Server responded with an error
	if err != nil {
		t.Error("1. An error was returned: ", err.Error())
	}
	// 2. Our email exists, If the server responded [false] which it shouldn't.. the test fails
	if !res.GetUsed() {
		t.Error("2. Username is used, should have returned true")
	}

	res, err = server.EmailUsed(context.Background(), &pb.EmailUsedRequest{Email: "incidrthreat_TEST@gmail.com"})
	// 3. Server responded with an error
	if err != nil {
		t.Error("3. An error was returned: ", err.Error())
	}
	// 3. Our email does not exist, If the server responded [true] which it shouldn't.. the test fails
	if res.GetUsed() {
		t.Error("3. User name is not used, should have returned false")
	}

}

func Test_authService_Register(t *testing.T) {
	// InMemroy MongoDB
	mongoServer, err := memongo.Start("4.0.5")
	assert.NoError(t, err)
	defer mongoServer.Stop()

	server := AuthServer{}

	insertTempUser(t, "incidrthreat", "incidrthreat@gmail.com", "incidrthreatpass", mongoServer.URIWithRandomDB())
	// Testing: Username in use
	_, err = server.Register(context.Background(), &pb.RegisterRequest{Email: "incidrthreatTEST@gmail.com", Username: "incidrthreat", Password: "incidrthreatpass"})
	if err.Error() != "Username in use" {
		t.Error("1. An error was returned")
	}
	// Testing: Email in use
	_, err = server.Register(context.Background(), &pb.RegisterRequest{Email: "incidrthreat@gmail.com", Username: "incidrTEST1", Password: "incidrthreatpass"})
	// We entered an email that was already in use, server should have responded with an error
	if err.Error() != "Email in use" {
		t.Error("2. An error was returned")
	}
	// Testing: Invalid Email
	_, err = server.Register(context.Background(), &pb.RegisterRequest{Email: "incidrthreat", Username: "incidrTEST2", Password: "incidrthreatpass"})
	// We entered an email that was already in use, server should have responded with an error
	if err.Error() != "That is not a valid email" {
		t.Error("3. An error was returned")
	}
	// Testing: Invalid Username
	_, err = server.Register(context.Background(), &pb.RegisterRequest{Email: "incidrthreatTEST3@gmail.com", Username: "inc", Password: "incidrthreatpass"})
	// We entered an invalid username, server should have responded with an error
	if err.Error() != "Username must be longer than 4 characters" {
		t.Error("4. An error was returned")
	}
	// Testing: Invalid Password
	_, err = server.Register(context.Background(), &pb.RegisterRequest{Email: "incidrthreatTEST4@gmail.com", Username: "incidrTEST4", Password: "incidr"})
	// We entered an invalid password, server should have responded with an error
	if err.Error() != "Password should be longer than 10 characters" {
		t.Error("5. An error was returned")
	}
	// Testing: Token Response
	res, err := server.Register(context.Background(), &pb.RegisterRequest{Email: "incidrthreatTEST5@gmail.com", Username: "incidrTEST5", Password: "incidr"})
	// We submitted a valid registration, server should have reported no errors
	if err != nil {
		t.Error("5. An error was returned")
	}
	// Check that the token is not empty
	if res.GetToken() != "" {
		// Check that the token matches our regex
		tokenRegEx, _ := regexp.MatchString("^[A-Za-z0-9-_=]+\\.[A-Za-z0-9-_=]+\\.?[A-Za-z0-9-_.+/=]*$", res.Token)
		if !tokenRegEx {
			t.Error("6. Invalid token returned")
		}
	} else {
		t.Error("7. No token returned")
	}

}
