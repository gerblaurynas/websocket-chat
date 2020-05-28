package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gerblaurynas/websocket-chat/internal/transport"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
)

func TestServer_ServeHTTP(t *testing.T) {
	srv := New(&zerolog.Logger{})
	s := httptest.NewServer(srv)

	// invalid route
	_ = dial(t, s, true, "test")

	conn := dial(t, s, false, "websocket?username=foo")
	require.NotNil(t, conn)
}

func TestServer_connectClient(t *testing.T) {
	srv := New(&zerolog.Logger{})
	s := httptest.NewServer(srv)

	_ = dial(t, s, false, "websocket?username=foo")

	// cannot upgrade connection
	resp, err := http.Get(s.URL + "/websocket?username=foobar")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUpgradeRequired, resp.StatusCode)

	// client already connected
	resp, err = http.Get(s.URL + "/websocket?username=foo")
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	_, ok := srv.connections["foo"]
	require.True(t, ok)
	require.Len(t, srv.connections, 1)
}

func TestServer_handleConnection(t *testing.T) {
	srv := New(&zerolog.Logger{})
	s := httptest.NewServer(srv)

	conn := dial(t, s, false, "websocket?username=foo")

	err := conn.Write(context.Background(), websocket.MessageText, []byte("foobar"))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	msgt, msg, err := conn.Read(ctx)
	cancel()
	assert.NoError(t, err)
	assert.Equal(t, websocket.MessageText, msgt)
	assert.Equal(t, `{"error": "invalid message"}`, string(msg))

	err = conn.Write(context.Background(), websocket.MessageText, []byte(`{"text": "hello"}`))
	require.NoError(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	beforeRead := time.Now()
	msgt, msg, err = conn.Read(ctx)
	cancel()
	assert.NoError(t, err)
	assert.Equal(t, websocket.MessageText, msgt)

	var received transport.Message
	err = json.Unmarshal(msg, &received)
	require.NoError(t, err)
	require.Equal(t, received.Username, "foo")
	require.Equal(t, received.Text, "hello")
	now := time.Now()
	require.True(t, received.Timestamp.After(beforeRead))
	require.True(t, received.Timestamp.Before(now))
}

func dial(t *testing.T, s *httptest.Server, fail bool, route string) *websocket.Conn {
	t.Helper()

	c, _, err := websocket.Dial(context.Background(), s.URL+"/"+route, nil)

	if fail {
		require.Error(t, err)
		return nil
	}

	require.NoError(t, err)

	return c
}
