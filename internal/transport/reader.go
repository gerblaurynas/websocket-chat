package transport

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"nhooyr.io/websocket"
)

type Reader struct{}

func (r *Reader) Read(conn *websocket.Conn, username string) (*Message, error) {
	messageType, msg, err := conn.Read(context.Background())
	if err != nil {
		return nil, err
	}

	if messageType != websocket.MessageText {
		_ = conn.Close(websocket.StatusNormalClosure, "non-text message received")

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
