package api

import (
	"fmt"
	"net/http"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", s.config.Welcome)

	// retrieve X-MS-CLIENT-PRINCIPAL-NAME header
	clientPrincipalName := r.Header.Get("X-Ms-Client-Principal-Name")
	if clientPrincipalName != "" {
		fmt.Fprintf(w, "Client Principal Name: %s\n", clientPrincipalName)
	}

}
