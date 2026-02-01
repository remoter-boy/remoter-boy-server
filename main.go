package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"remoter-boy-server/common"
	"remoter-boy-server/grpc"
	remoter "remoter-boy-server/proto_go"
	"syscall"
	"time"

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	sig := <-quit
	log.Printf("signal: %v", sig)

	done := make(chan struct{})
	go func() {
		srv.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Remoter-Boy Stopped")
	case <-time.After(10 * time.Second):
		log.Println("timeout! force stop")
		srv.Stop()
	default:
		err = grpc.DeleteClientAll(common.InitDatabase())
		if err != nil {
			log.Println("DeleteClientAll Error: " + err.Error())
		}
	}
}
