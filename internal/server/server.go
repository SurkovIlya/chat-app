package server

import (
	"context"
	"net/http"
	"time"

	chatserver "github.com/SurkovIlya/chat-app/internal/chat_server"
)

type SaveUser interface {
	SaveUser(userName string) error
}

type Server struct {
	httpServer *http.Server
	ChatServer *chatserver.ChatServer
	SaveUser   SaveUser
}

func New(port string, chs *chatserver.ChatServer, su SaveUser) *Server {
	s := &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    5000 * time.Millisecond,
			WriteTimeout:   5000 * time.Millisecond,
		},
		ChatServer: chs,
		SaveUser:   su,
	}

	s.initRoutes()

	return s
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// TODO: добавить хендлеры
func (s *Server) initRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", s.Connect)

	s.httpServer.Handler = mux
}
