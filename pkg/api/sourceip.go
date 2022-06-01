package api

import (
	"fmt"
	"net/http"
)

func (s *Server) sourceIpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", s.config.Welcome)
	fmt.Fprintf(w, "Source IP and port: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "X-Forwarded-For header: %s\n\n", r.Header.Get("X-Forwarded-For"))
	fmt.Fprintf(w, "All headers:\n\n")

	// print all HTTP headers
	for k, v := range r.Header {
		fmt.Fprintf(w, "HTTP header: %s: %s\n", k, v)
	}
}
