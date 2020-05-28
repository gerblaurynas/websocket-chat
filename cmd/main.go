package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gerblaurynas/websocket-chat/internal/server"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Stamp,
	}).With().Timestamp().Logger()

	if len(os.Args) != 2 {
		log.Fatal().Msg("port must be specified")
	}

	port := os.Args[1]

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatal().Msg("invalid port")
	}

	s := server.New(&log)
	err := http.ListenAndServe(":"+port, s)
	log.Fatal().Err(err).Msg("server terminated")
}
