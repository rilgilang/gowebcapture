package pkg

import (
	"context"
	socketio "github.com/googollee/go-socket.io"
)

type Socket interface {
	Broadcast(ctx context.Context, namespace string, message interface{})
}

type socketPkg struct {
	server *socketio.Server
}

func NewSocket(server *socketio.Server) Socket {
	return &socketPkg{
		server,
	}
}

func (s *socketPkg) Broadcast(ctx context.Context, namespace string, message interface{}) {
	s.server.BroadcastToNamespace(namespace, "broadcast", message)
}
