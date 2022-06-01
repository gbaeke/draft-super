package api

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type Pubsub struct {
	Pubsubname string `json:"pubsubname"`
	Topic      string `json:"topic"`
	Route      string `json:"route"`
}

func (s *Server) daprSubScribe(w http.ResponseWriter, r *http.Request) {
	// create JSON array as expected by Dapr
	pubsub := &[]Pubsub{{
		Pubsubname: s.config.Pubsub,
		Topic:      "mytopic",
		Route:      "/myroute",
	}}

	// return the JSON array to the caller (Dapr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pubsub)
	s.logger.Infow("Dapr called /dapr/subsribe route")

}

func (s *Server) myRoute(w http.ResponseWriter, r *http.Request) {
	// log the request body

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("Failed to read request body", zap.Error(err))
		return
	}
	s.logger.Infow("Received message", zap.ByteString("body", b))
	w.WriteHeader(http.StatusOK)
}
