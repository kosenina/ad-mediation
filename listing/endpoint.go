package listing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kosenina/ad-mediation/utils"
)

// Handler HTTP request handler
type Handler func(http.ResponseWriter, *http.Request)

// MakeGetAdNetworkListingEndpoint creates a handler for GET /adNetworkList requests
func MakeGetAdNetworkListingEndpoint(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse request parameter
		var documentID string
		queryValues := r.URL.Query()
		if len(queryValues) > 0 {
			dateString := queryValues.Get("date")
			layout := "2006-01-02"
			t, err := time.Parse(layout, dateString)

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Failed to parse provided date parameted, please use the right format: year-month-day (example: 2019-01-05)."}`))
				return
			}
			documentID = utils.GetAdNetworkListID(t)
		} else {
			documentID = utils.GetAdNetworkListID(time.Now())
		}

		// Get adNetworkList and return JSON object
		list, err := s.GetAdNetworkList(documentID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}
