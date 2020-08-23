package main

import (
	pb "github.com/golang-friends/slack-clone/protos"
)

type authServer {
}

func (as *authServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.AuthResponse, error){
	return pb.AuthResponse{}, nil
}
func (as *authServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return pb.AuthResponse{}, nil
}
func (as *authServer) UsernameUsed(ctx context.Context, in *pb.UsernameUsedRequest) (*pb.UsedResponse, error) {
	return pb.UsedResponse{}, nil
}
func (as *authServer) EmailUsed(ctx context.Context, in *pb.EmailUsedRequest) (*pbUsedResponse, error) {
	return pb.UsedResponse{}, nil	
}
