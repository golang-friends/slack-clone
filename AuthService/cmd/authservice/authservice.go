package authservice

import (
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/golang-friends/slack-clone/AuthService/models"
	pb "github.com/golang-friends/slack-clone/AuthService/protos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// AuthServer ...
type AuthServer struct {
}

// Login checks db for user info and returns a auth token
func (as *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.AuthResponse, error) {
	username, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	var user models.User
	models.Db.Collection("user").FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}).Decode(&user)
	if user == models.NilUser {
		return &pb.AuthResponse{}, errors.New("Wrong credentials provided")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &pb.AuthResponse{}, errors.New("Wrong credentials provided")
	}

	return &pb.AuthResponse{Token: user.GetToken()}, nil
}

// Register ...
func (as *AuthServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.AuthResponse, error) {
	username, password, email := in.GetUsername(), in.GetPassword(), in.GetEmail()
	emailRegEx, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if !emailRegEx {
		return &pb.AuthResponse{}, errors.New("That is not a valid email")
	}
	if len(username) < 4 {
		return &pb.AuthResponse{}, errors.New("Username must be longer than 4 characters")
	}
	if len(password) < 10 {
		return &pb.AuthResponse{}, errors.New("Password should be longer than 10 characters")
	}

	// Check if username is taken
	res, err := as.UsernameUsed(context.Background(), &pb.UsernameUsedRequest{Username: username})
	if err != nil {
		log.Println("Error returned", err.Error())
		return &pb.AuthResponse{}, errors.New("Something is wrong")
	}
	if res.GetUsed() {
		return &pb.AuthResponse{}, errors.New("Username in use")
	}

	// Check if email is taken
	res, err = as.EmailUsed(context.Background(), &pb.EmailUsedRequest{Email: email})
	if err != nil {
		log.Println("Error returned", err.Error())
		return &pb.AuthResponse{}, errors.New("Something is wrong")
	}
	if res.GetUsed() {
		return &pb.AuthResponse{}, errors.New("Email in use")
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := models.User{ID: primitive.NewObjectID(), Username: username, Email: email, Password: string(pw)}

	_, err = models.Db.Collection("user").InsertOne(context.Background(), newUser)
	if err != nil {
		log.Println("Error insterting new user: ", err.Error())
	}
	return &pb.AuthResponse{Token: newUser.GetToken()}, nil

}

// UsernameUsed ...
func (as *AuthServer) UsernameUsed(ctx context.Context, in *pb.UsernameUsedRequest) (*pb.UsedResponse, error) {
	username := in.GetUsername()
	var result models.User
	models.Db.Collection("user").FindOne(ctx, bson.M{"username": username}).Decode(&result)
	return &pb.UsedResponse{Used: result != models.NilUser}, nil
}

// EmailUsed ...
func (as *AuthServer) EmailUsed(ctx context.Context, in *pb.EmailUsedRequest) (*pb.UsedResponse, error) {
	email := in.GetEmail()
	var result models.User
	models.Db.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&result)
	return &pb.UsedResponse{Used: result != models.NilUser}, nil
}
