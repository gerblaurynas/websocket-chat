package transport

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gerblaurynas/websocket-chat/mocks"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
)

func TestReader_ReadEOF(t *testing.T) {
	reader := new(Reader)
	conn := new(mocks.Connection)

	conn.On("Read", context.Background()).
		Return(websocket.MessageBinary, []byte{}, errors.New("EOF")).
		Once()

	_, err := reader.Read(conn, "test")
	require.EqualError(t, err, "read error: EOF")
}

func TestReader_ReadInvalidType(t *testing.T) {
	reader := new(Reader)
	conn := new(mocks.Connection)

	conn.On("Read", context.Background()).
		Return(websocket.MessageBinary, []byte{}, nil).
		Once()

	_, err := reader.Read(conn, "test")
	require.EqualError(t, err, "non-text message received")
}

func TestReader_ReadInvalidMessage(t *testing.T) {
	reader := new(Reader)
	conn := new(mocks.Connection)

	conn.On("Read", context.Background()).
		Return(websocket.MessageText, []byte("foobar"), nil).
		Once()

	_, err := reader.Read(conn, "test")
	require.EqualError(t, err, "unmarshal message: invalid character 'o' in literal false (expecting 'a')")
}

func TestReader_Read(t *testing.T) {
	reader := new(Reader)
	conn := new(mocks.Connection)
	username := "foo"
	expectedMsg := "hello"

	conn.On("Read", context.Background()).
		Return(websocket.MessageText, []byte(fmt.Sprintf(`{"text":"%s"}`, expectedMsg)), nil).
		Once()

	msg, err := reader.Read(conn, username)

	require.NoError(t, err)
	require.Equal(t, expectedMsg, msg.Text)
	require.Equal(t, username, msg.Username)
	require.IsType(t, time.Time{}, msg.Timestamp)
}
