package grpc

import (
	remoter "remoter-boy-server/proto_go"
	"sync"
)

type Server struct {
	remoter.UnimplementedRemoterServer
	clients sync.Map
}
