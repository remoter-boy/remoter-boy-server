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
	log.Printf("[연결됨] 클라이언트: %s", client.ID)

	// 연결 종료 시 정리
	defer func() {
		s.clients.Delete(client.ID)
		log.Printf("[연결 끊김] 클라이언트: %s", client.ID)
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

		// 응답 전송
		if err := stream.Send(&remoter.NilResponseMsg{}); err != nil {
			return err
		}
	}
}
