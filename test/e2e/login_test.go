package e2e_test

import (
	"context"
	user "github.com/golang-friends/slack-clone/protos/user-svc"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"time"
)

const (
	localAddress = "127.0.0.1:9000"
)

func TestLogin(t *Testing.t) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err = user.Login(ctx, &user.LoginRequest{
		Username: 1,
		Email: 2,
		Password: 3
	})
	if err != nil {
		log.Fatalf("Login error: %v", err)
	}
	if !r.AuthResponse {
		log.Fatal("Failed to login")
	}
	log.Printf("Logged In: %t", r.AuthResponse)
}
