package adding

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kosenina/ad-mediation/models"
)

// Handler HTTP request handler
type Handler func(http.ResponseWriter, *http.Request)

// MakePutAdNetworkListingEndpoint creates a handler for PUT /adNetworkList requests
func MakePutAdNetworkListingEndpoint(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// get the body of our POST request
		// return the string response containing the request body
		reqBody, _ := ioutil.ReadAll(r.Body)
		var newAdNetworkList models.AdNetworkList
		var err = json.Unmarshal(reqBody, &newAdNetworkList)
		if err != nil {
			log.Println("ERROR: Failed to unmarshal request body.", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"errorMessage": "Invalid request body."}`))
			return
		}

		// Check if provided data is valid
		if newAdNetworkList.IsValid() == false {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"errorMessage": "List of AdNetwork items is not valid."}`))
			return
		}

		upsertErr := s.UpsertAdNetworkList(newAdNetworkList)
		if upsertErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"errorMessage": "Failed to update list of AdNetwork items."}`))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "Ad network list successfully updated."}`))
	}
}
