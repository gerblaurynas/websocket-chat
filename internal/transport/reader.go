package transport

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"nhooyr.io/websocket"
)

type Reader struct{}

type Connection interface {
	Read(ctx context.Context) (websocket.MessageType, []byte, error)
	Write(ctx context.Context, typ websocket.MessageType, p []byte) error
}

func (r *Reader) Read(conn Connection, username string) (*Message, error) {
	messageType, msg, err := conn.Read(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "read error")
	}

	if messageType != websocket.MessageText {
		return nil, UserError{
			error:       errors.New("non-text message received"),
			userMessage: "invalid message",
		}
	}

	var m Message

	err = json.Unmarshal(msg, &m)
	if err != nil {
		return nil, UserError{
			error:       errors.Wrap(err, "unmarshal message"),
			userMessage: "invalid message",
		}
	}

	m.Username = username
	m.Timestamp = time.Now()

	return &m, nil
}
