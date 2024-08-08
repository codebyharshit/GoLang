package db

import (
	"database/sql"

	"github.com/codebyharshit/real-time-analytics/internal/core/entities"
)

type DatabaseStorage struct {
	DB *sql.DB
}

func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{DB: db}
}

func (s *DatabaseStorage) SaveTrade(trade entities.Trade) error {
	_, err := s.DB.Exec(
		"INSERT INTO trades (id, timestamp, symbol, quantity, price, side) VALUES ($1, $2, $3, $4, $5, $6)",
		trade.ID, trade.Timestamp, trade.Symbol, trade.Quantity, trade.Price, trade.Side,
	)
	if err != nil {
		return err
	}

	// Update portfolio
	return s.updatePortfolio(trade)
}

func (s *DatabaseStorage) updatePortfolio(trade entities.Trade) error {
	var portfolio entities.Portfolio

	// Fetch existing portfolio
	err := s.DB.QueryRow("SELECT id, cash, total_value FROM portfolio WHERE id = $1", trade.ID).Scan(
		&portfolio.ID, &portfolio.Cash, &portfolio.TotalValue,
	)
	if err != nil {
		return err
	}

	// Update holdings and cash based on trade side
	if trade.Side == "buy" {
		portfolio.Holdings[trade.Symbol] += trade.Quantity
		portfolio.Cash -= trade.Quantity * trade.Price
	} else {
		portfolio.Holdings[trade.Symbol] -= trade.Quantity
		portfolio.Cash += trade.Quantity * trade.Price
	}

	portfolio.TotalValue = portfolio.Cash
	for symbol, quantity := range portfolio.Holdings {
		var price float64
		err := s.DB.QueryRow("SELECT price FROM trades WHERE symbol = $1 ORDER BY timestamp DESC LIMIT 1", symbol).Scan(&price)
		if err != nil {
			return err
		}
		portfolio.TotalValue += quantity * price
	}

	// Save updated portfolio
	_, err = s.DB.Exec("UPDATE portfolio SET cash = $1, total_value = $2 WHERE id = $3",
		portfolio.Cash, portfolio.TotalValue, portfolio.ID,
	)
	return err
}

func (s *DatabaseStorage) GetPortfolio() (entities.Portfolio, error) {
	var portfolio entities.Portfolio

	rows, err := s.DB.Query("SELECT symbol, quantity FROM holdings WHERE portfolio_id = $1", portfolio.ID)
	if err != nil {
		return portfolio, err
	}
	defer rows.Close()

	portfolio.Holdings = make(map[string]float64)
	for rows.Next() {
		var symbol string
		var quantity float64
		if err := rows.Scan(&symbol, &quantity); err != nil {
			return portfolio, err
		}
		portfolio.Holdings[symbol] = quantity
	}

	return portfolio, nil
}
