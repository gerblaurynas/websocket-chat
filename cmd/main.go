package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gerblaurynas/websocket-chat/internal/server"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Stamp,
	}).With().Timestamp().Logger()

	s := server.New(&log)
	log.Info().Msg("server is running...")

	err := http.ListenAndServe(":80", s)
	log.Fatal().Err(err).Msg("server terminated")
}
