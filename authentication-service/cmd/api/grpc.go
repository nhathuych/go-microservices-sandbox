package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/nhathuych/go-microservices-sandbox/authentication-service/auth"
	"github.com/nhathuych/go-microservices-sandbox/authentication-service/proto"
	"google.golang.org/grpc"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer
	model auth.Models
}

func (us *UserServer) GetUserByID(ctx context.Context, req *proto.UserRequest) (*proto.UserResponse, error) {
	userID := req.GetId()

	user, err := us.model.User.GetOne(int(userID))
	if err != nil {
		log.Println(err)
		response := proto.UserResponse{User: nil}
		return &response, err
	}

	response := proto.UserResponse{User: &proto.User{
		Id:         int64(user.ID),
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		UserActive: int32(user.Active),
	}}

	return &response, nil
}

func (app *Config) gRPCListen() {
	gRpcPort := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()

	userServer := &UserServer{model: app.Models}
	proto.RegisterUserServiceServer(srv, userServer)
	log.Printf("gRPC Server started on port %d", gRpcPort)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
