package api

import (
	"io"
	"net/http"

	"go.uber.org/zap"
)

// implements /mqtt route to accept MQTT events via Dapr
func (s *Server) mqtt(w http.ResponseWriter, r *http.Request) {
	// log the request body
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("Failed to read request body", zap.Error(err))
		return
	}
	s.logger.Infow("Received MQTT message", zap.ByteString("body", b))
	w.WriteHeader(http.StatusOK)
}
