package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gerblaurynas/websocket-chat/internal/transport"
	"github.com/rs/zerolog"
	"nhooyr.io/websocket"
)

type Server struct {
	log         zerolog.Logger
	connections map[string]*websocket.Conn
	reader      *transport.Reader
	writer      *transport.Writer
}

func New(log *zerolog.Logger) *Server {
	return &Server{
		log:         *log,
		connections: make(map[string]*websocket.Conn),
		reader:      &transport.Reader{},
		writer:      &transport.Writer{},
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/websocket" {
		s.handleConnection(w, r)
		return
	}

	http.NotFoundHandler().ServeHTTP(w, r)
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	username, c, err := s.connectClient(w, r)
	if err != nil {
		s.log.Warn().Err(err).Msg(err.Error())
		return
	}

	defer s.disconnectClient(username, c)

	for {
		input, err := s.reader.Read(c)
		if err != nil {
			s.log.Warn().Msg(err.Error())
			return
		}

		err = s.writer.Send(input, username, s.connections)
		if err != nil {
			s.log.Warn().Msg(err.Error())
			return
		}
	}
}

func (s *Server) connectClient(w http.ResponseWriter, r *http.Request) (username string, c *websocket.Conn, err error) {
	username = r.URL.Query().Get("username")
	if username == "" {
		err = errors.New("username not provided")
		return
	}

	if _, exists := s.connections[username]; exists {
		err = errors.New("client with this username already connected")
		return
	}

	// allow connections from chrome extensions
	options := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}

	c, err = websocket.Accept(w, r, options)
	if err != nil {
		err = errors.New("cannot upgrade connection")
		return
	}

	s.connections[username] = c
	s.log.Info().Msg(fmt.Sprintf("%s connected", username))

	return
}

func (s *Server) disconnectClient(username string, c *websocket.Conn) {
	_ = c.Close(websocket.StatusInternalError, "internal error")
	delete(s.connections, username)
	s.log.Info().Msg(fmt.Sprintf("%s disconnected", username))
}
