package server

import (
	"fmt"
	"net/http"

	"github.com/gerblaurynas/websocket-chat/internal/transport"
	"github.com/pkg/errors"
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
			switch err.(type) {
			case transport.UserError:
				s.writer.SendError(err, c)
				continue
			default:
				s.log.Warn().Msg(err.Error())
				return
			}
		}

		err = s.writer.SendAll(input, username, s.connections)
		if err != nil {
			s.log.Warn().Msg(err.Error())
			return
		}
	}
}

func (s *Server) connectClient(w http.ResponseWriter, r *http.Request) (username string, c *websocket.Conn, err error) {
	// allow connections from chrome extensions
	options := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}

	c, err = websocket.Accept(w, r, options)
	if err != nil {
		err = errors.Wrap(err, "cannot upgrade connection")
		return
	}

	username = r.URL.Query().Get("username")
	if username == "" {
		err = errors.New("username not provided")
		s.abortWithError(err, c)

		return
	}

	if _, exists := s.connections[username]; exists {
		err = errors.New(fmt.Sprintf("client with username %s is already connected", username))
		s.abortWithError(err, c)

		return
	}

	s.connections[username] = c
	s.log.Info().Msg(fmt.Sprintf("%s connected", username))

	return username, c, nil
}

func (s *Server) disconnectClient(username string, c *websocket.Conn) {
	_ = c.Close(websocket.StatusInternalError, "internal error")

	delete(s.connections, username)
	s.log.Info().Msg(fmt.Sprintf("%s disconnected", username))
}

func (s *Server) abortWithError(reason error, c *websocket.Conn) {
	s.writer.SendError(reason, c)
	_ = c.Close(websocket.StatusInternalError, reason.Error())
}
