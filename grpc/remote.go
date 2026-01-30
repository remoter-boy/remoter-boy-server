package grpc

import (
	"database/sql"
	"io"
	"log"
	"remoter-boy-server/common"
	remoter "remoter-boy-server/proto_go"
	"time"
)

type ClientInfo struct {
	ID       string
	Stream   remoter.Remoter_ConnectServer
	LastPing time.Time
}

func CreateClient(db *sql.DB, clientId string) (*sql.Result, error) {
	query := `INSERT INTO public.tb_client ("client_id") VALUES ($1)`
	result, err := db.Exec(query, clientId)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func DeleteClient(db *sql.DB, clientId string) (*sql.Result, error) {
	query := "DELETE FROM public.tb_client where client_id = $1"
	result, err := db.Exec(query, clientId)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Server) Connect(stream remoter.Remoter_ConnectServer) error {
	db := common.InitDatabase()

	if db == nil {
		panic("Database Connect Error")
	}

	defer db.Close()

	recv, err := stream.Recv()

	if err != nil {
		log.Println("stream.Recv() Error:" + err.Error())
		return err
	}

	_, err = CreateClient(db, recv.ClientId)

	if err != nil {
		log.Println("Client Insert Error:" + err.Error())
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
		_, err = DeleteClient(db, client.ID)

		if err != nil {
			log.Println("Client Delete Error: " + err.Error())
		}
		log.Printf("[Disconnet] Client: %s", client.ID)
	}()

	for {
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		client.LastPing = time.Now()

		// send Response
		if err := stream.Send(&remoter.NilResponseMsg{}); err != nil {
			return err
		}
	}
}
