package listing

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kosenina/ad-mediation/utils"
)

// Handler HTTP request handler
type Handler func(c *gin.Context)

// MakeGetAdNetworkListingEndpoint creates a handler for GET /adNetworkList requests
func MakeGetAdNetworkListingEndpoint(s Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Parse request parameter
		var documentID string
		queryDate := c.DefaultQuery("date", "")
		if queryDate == "" {
			documentID = utils.GetAdNetworkListID(time.Now())
		} else {
			layout := "2006-01-02"
			t, err := time.Parse(layout, queryDate)

			if err != nil {
				log.Printf(fmt.Sprintf("ERROR: Failed to parse provided query date parameter: %s", queryDate))
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			documentID = utils.GetAdNetworkListID(t)
		}

		// Get adNetworkList and return JSON object
		list, err := s.GetAdNetworkList(documentID)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, &list)
	}
}
