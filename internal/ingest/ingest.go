package ingest

import (
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Start(pErr chan<- error) {
	pErr <- startWebServer()
}

func startWebServer() error {
	server := http.NewServeMux()
	server.HandleFunc("/new_rtc_session", newRTCSessionHandler)
	addr := "0.0.0.0:8899"
	log.Info().Msg("Start web server on " + addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		return err
	}
	return errors.New("unxpected exit from ingest http server")
}
