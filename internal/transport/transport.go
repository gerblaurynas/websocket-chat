package transport

import (
	"context"
	"encoding/json"
	"errors"

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
		return nil, errors.New("non-text message received")
	}

	var input Input

	err = json.Unmarshal(msg, &input)
	if err != nil {
		return nil, errors.New("invalid message format")
	}

	return &input, nil
}
