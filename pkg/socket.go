package pkg

import (
	"context"
	socketio "github.com/googollee/go-socket.io"
)

type Socket interface {
	VideoProcessingComplete(ctx context.Context, namespace string, message interface{})
	VideoProcessingFail(ctx context.Context, namespace string, message interface{})
}

type socketPkg struct {
	server *socketio.Server
}

func NewSocket(server *socketio.Server) Socket {
	return &socketPkg{
		server,
	}
}

func (s *socketPkg) VideoProcessingComplete(ctx context.Context, namespace string, message interface{}) {
	s.server.BroadcastToNamespace(namespace, "video-complete", message)
}

func (s *socketPkg) VideoProcessingFail(ctx context.Context, namespace string, message interface{}) {
	s.server.BroadcastToNamespace(namespace, "video-error", message)
}
