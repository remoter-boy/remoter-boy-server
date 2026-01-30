package main

import (
	"log"
	"net"
	"remoter-boy-server/common"
	"remoter-boy-server/grpc"
	remoter "remoter-boy-server/proto_go"

	grpcs "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := common.GetConfig()
	lis, err := net.Listen("tcp", ":"+config.Port)

	if err != nil {
		panic(err)
	}

	srv := grpcs.NewServer()

	remoter.RegisterRemoterServer(srv, &grpc.Server{})
	reflection.Register(srv)
	log.Println("Remoter-Boy Start")

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
