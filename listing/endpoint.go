package listing

import (
	"encoding/json"
	"net/http"
)

// Handler HTTP request handler
type Handler func(http.ResponseWriter, *http.Request)

// MakeGetAdNetworkListingEndpoint creates a handler for GET /adNetworkList requests
func MakeGetAdNetworkListingEndpoint(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		list, err := s.GetAdNetworkList()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Failed to get ad network list."}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}
