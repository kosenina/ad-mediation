package listing

import "github.com/kosenina/ad-mediation/models"

// Service provides ad network listing operations
type Service interface {
	GetAdNetworkList() (models.AdNetworkList, error)
}

type service struct {
	repo models.Repository
}

// NewService creates service with the necessary dependencies
func NewService(repo models.Repository) Service {
	return &service{repo}
}

// GetAdNetworkList returns AdNetworkList struct
func (s *service) GetAdNetworkList() (models.AdNetworkList, error) {
	return s.repo.Get()
}
