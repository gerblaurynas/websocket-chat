package transport

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"nhooyr.io/websocket"
)

type Reader struct{}

func (r *Reader) Read(conn *websocket.Conn) (*Input, error) {
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

	var input Input

	err = json.Unmarshal(msg, &input)
	if err != nil {
		return nil, UserError{
			error:       errors.Wrap(err, "unmarshal message"),
			userMessage: "invalid message",
		}
	}

	return &input, nil
}
