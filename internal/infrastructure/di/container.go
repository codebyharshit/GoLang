package di

import (
	"github.com/codebyharshit/real-time-analytics/internal/app/analytics"
	"github.com/codebyharshit/real-time-analytics/internal/app/risk"
	"github.com/codebyharshit/real-time-analytics/internal/app/trading"
	"github.com/codebyharshit/real-time-analytics/internal/infrastructure/db"
	"github.com/codebyharshit/real-time-analytics/pkg/config"
)

type Container struct {
	TraderService      *trading.TraderService
	RiskManagerService *risk.RiskManagerService
	PredictorService   *analytics.PredictorService
}

func NewContainer() (*Container, error) {
	dbConn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}

	storage := db.NewDatabaseStorage(dbConn)
	riskMgr := risk.NewRiskManagerService(storage)
	predictor := analytics.NewPredictorService("http://localhost:5000/predict")
	traderService := trading.NewTraderService(storage, riskMgr, predictor)

	return &Container{
		TraderService:      traderService,
		RiskManagerService: riskMgr,
		PredictorService:   predictor,
	}, nil
}
