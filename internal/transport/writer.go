package transport

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"nhooyr.io/websocket"
)

type Writer struct{}

func (w *Writer) Send(i *Input, author string, connections map[string]*websocket.Conn) error {
	output := &Output{
		Input:     i,
		Timestamp: time.Now(),
		Username:  author,
	}

	outMsg, err := json.Marshal(output)
	if err != nil {
		return errors.New("failed to build message")
	}

	for i := range connections {
		err = connections[i].Write(context.Background(), websocket.MessageText, outMsg)
		if err != nil {
			return err
		}
	}

	return nil
}
