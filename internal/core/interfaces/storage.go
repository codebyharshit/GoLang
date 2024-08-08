package interfaces

import "github.com/codebyharshit/real-time-analytics/internal/core/entities"

type Storage interface {
	SaveTrade(trade entities.Trade) error
	GetPortfolio() (entities.Portfolio, error)
}
