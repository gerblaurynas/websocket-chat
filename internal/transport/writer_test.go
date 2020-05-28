package transport

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gerblaurynas/websocket-chat/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
)

func TestWriter_SendAll(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2020-05-28T15:00:00Z")
	require.NoError(t, err)
	msg := &Message{
		Text:      "hello",
		Timestamp: timestamp,
		Username:  "foobar",
	}
	rawMsg := []byte(`{"text":"hello","timestamp":"2020-05-28T15:00:00Z","username":"foobar"}`)

	writer := new(Writer)
	conn1 := new(mocks.Connection)
	conn1.On("Write", context.Background(), websocket.MessageText, rawMsg).
		Return(nil).
		Once()

	conn2 := new(mocks.Connection)
	conn2.On("Write", context.Background(), websocket.MessageText, rawMsg).
		Return(nil).
		Once()

	connections := make(map[string]Connection)
	connections["foobar"] = conn1
	connections["baz"] = conn2

	err = writer.SendAll(msg, connections)
	require.NoError(t, err)
}

func TestWriter_SendError(t *testing.T) {
	writer := new(Writer)
	conn := new(mocks.Connection)
	expected := `{"error": "test error"}`

	conn.On("Write", context.Background(), websocket.MessageText, mock.MatchedBy(func(msg []byte) bool {
		return string(msg) == expected
	})).
		Return(nil).
		Once()

	err := errors.New("test error")

	writer.SendError(err, conn)
}
func TestWriter_SendUserError(t *testing.T) {
	writer := new(Writer)
	conn := new(mocks.Connection)
	expected := `{"error": "test user facing error"}`

	conn.On("Write", context.Background(), websocket.MessageText, mock.MatchedBy(func(msg []byte) bool {
		return string(msg) == expected
	})).
		Return(nil).
		Once()

	err := UserError{
		error:       errors.New("test error"),
		userMessage: "test user facing error",
	}

	writer.SendError(err, conn)
}
