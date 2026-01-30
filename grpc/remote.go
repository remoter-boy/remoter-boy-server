package grpc

import (
	"io"
	"log"
	remoter "remoter-boy-server/proto_go"
	"time"
)

type ClientInfo struct {
	ID       string
	Stream   remoter.Remoter_ConnectServer
	LastPing time.Time
}

func (s *Server) Connect(stream remoter.Remoter_ConnectServer) error {
	recv, err := stream.Recv()

	if err != nil {
		return err
	}

	client := &ClientInfo{
		ID:       recv.ClientId,
		Stream:   stream,
		LastPing: time.Now(),
	}
	s.clients.Store(client.ID, client)
	log.Printf("[Connect] Client: %s", client.ID)

	// Clean up at the end of the connection
	defer func() {
		s.clients.Delete(client.ID)
		log.Printf("[Disconnet] Client: %s", client.ID)
	}()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Println(msg.ClientId)

		client.LastPing = time.Now()

		// send Response
		if err := stream.Send(&remoter.NilResponseMsg{}); err != nil {
			return err
		}
	}
}
