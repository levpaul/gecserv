package ingest

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

var pipeErr chan error

func Start(pErr chan<- error) error {
	go startWebServer()
	return nil
}

func startWebServer() {
	server := http.NewServeMux()
	server.HandleFunc("/new_rtc_session", newRTCSessionHandler)
	addr := "127.0.0.1:8899"
	log.Info().Msg("Start web server on " + addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		pipeErr <- err
	}
}
