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
// @Summary Get ad network list
// @Description get ad network list for queried time, current (for today) ad network list is returned if parameter is not provided
// @Tags adNetworkList
// @Accept  json
// @Produce  json
// @Param date query string false "provide time to fetch list of ad newtworks for specified time" Format(date)
// @Success 200 {object} models.AdNetworkList
// @Router /adNetworkList [get]
func MakeGetAdNetworkListingEndpoint(s Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

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
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"message": "Failed to parse provided date parameted, please use the right format: year-month-day (example: 2019-01-05)."})
				return
			}
			documentID = utils.GetAdNetworkListID(t)
		}

		// Get adNetworkList and return JSON object
		list, err := s.GetAdNetworkList(documentID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Requested document does not exists."})
			return
		}

		c.JSON(http.StatusOK, &list)
	}
}
