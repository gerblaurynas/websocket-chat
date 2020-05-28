package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"nhooyr.io/websocket"
)

type Writer struct{}

func (w *Writer) SendAll(m *Message, connections map[string]*websocket.Conn) error {
	outMsg, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "failed to build message")
	}

	for i := range connections {
		err = connections[i].Write(context.Background(), websocket.MessageText, outMsg)
		if err != nil {
			return errors.Wrap(err, "send for connections")
		}
	}

	return nil
}

func (w *Writer) SendError(msg error, c *websocket.Conn) {
	var text string

	userError, ok := msg.(UserError)
	if ok {
		text = userError.userMessage
	} else {
		text = msg.Error()
	}

	message := fmt.Sprintf(`{"error": "%s"}`, text)
	_ = c.Write(context.Background(), websocket.MessageText, []byte(message))
}
