package main

import (
	"context"
	"log"
	"net"

	pb "github.com/golang-friends/slack-clone/protos"
	"google.golang.org/grpc"
)

type authServer struct {
}

func (as *authServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{}, nil
}
func (as *authServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{}, nil
}
func (as *authServer) UsernameUsed(ctx context.Context, in *pb.UsernameUsedRequest) (*pb.UsedResponse, error) {
	return &pb.UsedResponse{}, nil
}
func (as *authServer) EmailUsed(ctx context.Context, in *pb.EmailUsedRequest) (*pb.UsedResponse, error) {
	return &pb.UsedResponse{}, nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, &authServer{})

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Error creating listener: ", err.Error())
	}

	server.Serve(listener)
}
