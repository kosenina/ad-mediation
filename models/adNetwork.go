package models

// AdNetwork metadata
type AdNetwork struct {
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

// AdNetworkList list of AdNetworks
type AdNetworkList struct {
	Items []AdNetwork `json:"items"`
}

// IsValid check wether object data is correct
func (a AdNetworkList) IsValid() bool {
	// At least one AdNetwork item must be provided
	itemsLen := len(a.Items)
	if itemsLen < 1 {
		return false
	}

	// AdNetwork items ranks must be unique and in correct range
	rangeSum := 0
	uniqueRanks := make(map[int]int)
	for _, s := range a.Items {
		uniqueRanks[s.Rank] = 1
		rangeSum += s.Rank
		if len(s.Name) < 1 {
			return false
		}
	}

	// Check if ranks are in correct sequence
	if rangeSum != (itemsLen*(itemsLen-1))/2 {
		return false
	}

	// Check that ranks are not duplicated
	if len(uniqueRanks) != itemsLen {
		return false
	}
	return true
}

// Repository provides access to the ad network storage.
type Repository interface {
	// Get returns ad network saved in DB.
	Get() (AdNetworkList, error)
	// Add saves an ad network into the DB.
	Upsert(AdNetworkList) error
}
