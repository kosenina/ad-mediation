package utils

import (
	"fmt"
	"time"
)

// GetAdNetworkListID return deterministic ID for the adNetworkList document
func GetAdNetworkListID(date time.Time) string {
	year, month, day := date.Date()
	s := fmt.Sprintf("adNetList_%d:%d:%d", day, month, year)
	return s
}
