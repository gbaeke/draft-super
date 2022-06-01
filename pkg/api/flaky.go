package api

import (
	"math/rand"
	"net/http"
)

func (s *Server) flakyHandler(w http.ResponseWriter, r *http.Request) {
	// return status code 500 10% of the time
	if rand.Intn(10) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
