package adding

import "github.com/kosenina/ad-mediation/models"

// Service provides ad network listing operations
type Service interface {
	UpsertAdNetworkList(models.AdNetworkList) error
}

type service struct {
	repo models.Repository
}

// NewService creates service with the necessary dependencies
func NewService(repo models.Repository) Service {
	return &service{repo}
}

// GetAdNetworkList returns AdNetworkList struct
func (s *service) UpsertAdNetworkList(data models.AdNetworkList) error {
	return s.repo.Upsert(data)
}
