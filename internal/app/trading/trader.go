package trading

import (
	"github.com/codebyharshit/real-time-analytics/internal/app/analytics"
	"github.com/codebyharshit/real-time-analytics/internal/core/entities"
	"github.com/codebyharshit/real-time-analytics/internal/core/interfaces"
)

type TraderService struct {
	storage   interfaces.Storage
	riskMgr   interfaces.RiskManager
	predictor *analytics.PredictorService
}

func NewTraderService(storage interfaces.Storage, riskMgr interfaces.RiskManager, predictor *analytics.PredictorService) *TraderService {
	return &TraderService{storage: storage, riskMgr: riskMgr, predictor: predictor}
}

func (s *TraderService) ExecuteTrade(trade entities.Trade) error {
	if err := s.riskMgr.EvaluateRisk(entities.MarketData{
		Symbol: trade.Symbol,
		Price:  trade.Price,
	}); err != nil {
		return err
	}

	// Get the prediction for the trade
	prediction, err := s.predictor.Predict(entities.MarketData{
		SMA_50:  trade.Price * 1.05, // Example placeholder calculation
		SMA_200: trade.Price * 0.95, // Example placeholder calculation
	})
	if err != nil {
		return err
	}

	// Logic based on prediction
	if prediction == 1 {
		// Execute trade if prediction is positive
		return s.storage.SaveTrade(trade)
	}

	return nil // Do not execute trade if prediction is not positive
}

func (s *TraderService) GetPortfolio() (entities.Portfolio, error) {
	return s.storage.GetPortfolio()
}
