// internal/infrastructure/storage/in_memory.go
package storage

import (
	"log"

	"github.com/codebyharshit/real-time-analytics/internal/core/entities"
)

type InMemoryStorage struct {
	trades    []entities.Trade
	portfolio entities.Portfolio
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		trades:    make([]entities.Trade, 0),
		portfolio: entities.Portfolio{Holdings: make(map[string]float64)},
	}
}

func (s *InMemoryStorage) SaveTrade(trade entities.Trade) error {
	log.Println("Adding trade to in-memory storage")
	s.trades = append(s.trades, trade)

	// Update portfolio (simplified logic)
	if trade.Side == "buy" {
		s.portfolio.Holdings[trade.Symbol] += trade.Quantity
		s.portfolio.Cash -= trade.Quantity * trade.Price
	} else {
		s.portfolio.Holdings[trade.Symbol] -= trade.Quantity
		s.portfolio.Cash += trade.Quantity * trade.Price
	}

	// Recalculate total value of the portfolio
	s.portfolio.TotalValue = s.portfolio.Cash
	for symbol, quantity := range s.portfolio.Holdings {
		for _, t := range s.trades {
			if t.Symbol == symbol {
				s.portfolio.TotalValue += quantity * t.Price
				break
			}
		}
	}
	log.Println("Trade added to storage successfully")
	return nil
}

func (s *InMemoryStorage) GetPortfolio() (entities.Portfolio, error) {
	log.Println("Fetching portfolio from storage")
	return s.portfolio, nil
}
