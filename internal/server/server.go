package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gerblaurynas/websocket-chat/internal/transport"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"nhooyr.io/websocket"
)

type Server struct {
	log         zerolog.Logger
	connections map[string]transport.Connection
	reader      *transport.Reader
	writer      *transport.Writer
	lock        sync.Locker
}

func New(log *zerolog.Logger) *Server {
	return &Server{
		log:         *log,
		connections: make(map[string]transport.Connection),
		reader:      new(transport.Reader),
		writer:      new(transport.Writer),
		lock:        new(sync.Mutex),
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
		s.log.Warn().Msg(err.Error())
		return
	}

	defer s.disconnectClient(username, c)

	for {
		msg, err := s.reader.Read(c, username)
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

		s.lock.Lock()
		err = s.writer.SendAll(msg, s.connections)
		s.lock.Unlock()
		if err != nil {
			s.log.Warn().Msg(err.Error())
			return
		}
	}
}

func (s *Server) connectClient(w http.ResponseWriter, r *http.Request) (username string, c *websocket.Conn, err error) {
	username = r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username not found", 400)
		err = errors.New("username not found")
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.connections[username]; ok {
		http.Error(w, "username already connected", 409)
		err = errors.New("username already connected")
		return
	}

	// allow connections from chrome extensions
	options := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}

	c, err = websocket.Accept(w, r, options)
	if err != nil {
		err = errors.Wrap(err, "cannot upgrade connection")
		return
	}

	s.connections[username] = c
	s.log.Info().Msg(fmt.Sprintf("%s connected", username))

	return username, c, nil
}

func (s *Server) disconnectClient(username string, c *websocket.Conn) {
	_ = c.Close(websocket.StatusInternalError, "internal error")

	s.lock.Lock()
	delete(s.connections, username)
	s.lock.Unlock()
	s.log.Info().Msg(fmt.Sprintf("%s disconnected", username))
}
