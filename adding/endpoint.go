package adding

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kosenina/ad-mediation/models"
)

// Handler HTTP request handler
type Handler func(c *gin.Context)

// MakePutAdNetworkListingEndpoint creates a handler for PUT /adNetworkList requests
func MakePutAdNetworkListingEndpoint(s Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// get the body of our POST request
		// return the string response containing the request body
		reqBody, bodyErr := c.GetRawData()
		if bodyErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ad network list must be provided in request body."})
			return
		}
		var newAdNetworkList models.AdNetworkList
		var err = json.Unmarshal(reqBody, &newAdNetworkList)
		if err != nil {
			log.Println("ERROR: Failed to unmarshal request body.", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Request body cannot be parsed."})
			return
		}

		// Check if provided data is valid
		if newAdNetworkList.IsValid() == false {
			log.Println("ERROR: PUT-ed ad newtwork list is not valid")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ad network list is not valid."})
			return
		}

		upsertErr := s.UpsertAdNetworkList(newAdNetworkList)
		if upsertErr != nil {
			log.Println("ERROR: Failed to update list of ad networks")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to upsert ad network list."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Ad network list successfully updated."})
	}
}
