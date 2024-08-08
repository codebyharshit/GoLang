package risk

import (
	"github.com/codebyharshit/real-time-analytics/internal/core/entities"
	"github.com/codebyharshit/real-time-analytics/internal/core/interfaces"
)

type RiskManagerService struct {
	storage interfaces.Storage
}

func NewRiskManagerService(storage interfaces.Storage) *RiskManagerService {
	return &RiskManagerService{storage: storage}
}

func (s *RiskManagerService) EvaluateRisk(data entities.MarketData) error {
	return nil
}
