package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type State struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

// @Summary Save state
// @Description Save state to configured state store
// @Accept json
// @Produces json
// @Router /state [post]
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "Error reading or unmarshalling request body"
// @Failure 500 {string} string "Error writing to statestore"
func (s *Server) saveState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.logger.Infow("Method not allowed on saveState",
			zap.String("error", r.Method))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Infow("Could not save state: error reading request body",
			zap.Error(err))
		return
	}

	// unmarshal the request body, expecting a JSON object with a key and data
	var state State
	err = json.Unmarshal(body, &state)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Infow("Could not save state: invalid request body",
			zap.Error(err))
		return
	}

	// return error if key or data is empty
	if state.Key == "" || state.Data == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Infow("Could not save state: key or data is empty",
			zap.String("key", state.Key),
			zap.String("data", state.Data))
		return
	}

	// write data to Dapr statestore
	ctx := r.Context()
	if err := s.daprClient.SaveState(ctx, s.config.Statestore, state.Key, []byte(state.Data)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Infow("Could not write to statestore", zap.String("key", state.Key),
			zap.String("statestore", s.config.Statestore))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		s.logger.Infow("Successfully wrote to statestore", zap.String("key", state.Key))
	}

}

// @Summary Read state
// @Description Read state from configured state store
// @Accept json
// @Produces json
// @Router /state [get]
func (s *Server) readState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.logger.Infow("Method not allowed on readState", zap.String("method", r.Method))
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Infow("Could not read from statestore: key missing from request")
		return
	}

	// read data from Dapr statestore
	ctx := r.Context()
	data, err := s.daprClient.GetState(ctx, s.config.Statestore, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Infow("Could not read from statestore", "key", key, zap.Error(err))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(data.Value)
		s.logger.Infow("Successfully read from statestore",
			zap.String("key", key))
	}

}
