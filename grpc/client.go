package grpc

import (
	"context"
	"database/sql"
	remoter "remoter-boy-server/proto_go"
	"time"
)

type ClientInfo struct {
	ID       string
	Stream   remoter.Remoter_ConnectServer
	LastPing time.Time
	Cancel   context.CancelFunc
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
