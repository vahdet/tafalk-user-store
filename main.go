package main

import (
	"fmt"
	"github.com/vahdet/go-user-store/app"
	"github.com/vahdet/go-user-store/dal"
	"github.com/vahdet/go-user-store/services"
	"net"
	pb "github.com/vahdet/go-user-store/proto"
	"google.golang.org/grpc"
	"github.com/vahdet/go-user-store/grpcserver"
	"google.golang.org/grpc/reflection"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Deprecated read from env variables instead of files
	/*
	if err := app.LoadConfig("./resources"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}
	if err := app.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(fmt.Errorf("failed to read the error message file: %s", err))
	}
	*/

	// Initialize DataStore
	dal.InitDataStoreClient()
	defer dal.CloseDataStoreClient()

	// Serve GRPC
	lis, err := net.Listen("tcp", app.Config.GrpcPort)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &grpcserver.UserServer{Service: services.NewUserService(dal.NewUserDal())})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
