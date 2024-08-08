package interfaces

import (
	"github.com/codebyharshit/real-time-analytics/internal/core/entities"
)

type MarketDataProcessor interface {
	Process(data entities.MarketData) error
}
type Trader interface {
	ExecuteTrade(trade entities.Trade) error
	GetPortfolio() (entities.Portfolio, error)
}
type RiskManager interface {
	EvaluateRisk(data entities.MarketData) error
}
