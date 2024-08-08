# Algorithmic Trading System

## Overview

This project is an algorithmic trading system built using Go and Python. It fetches historical market data, trains a machine learning model for predictions, and integrates these predictions into a Go-based trading system. The system makes rapid trading decisions based on market data and performs risk management by monitoring market conditions and portfolio performance in real-time.

## Features

- **Data Ingestion**: Fetches historical market data from the Polygon API.
- **Model Training**: Trains a machine learning model to predict trading signals based on market data.
- **Prediction Service**: A Flask-based API that serves model predictions.
- **Trading System**: A Go application that integrates predictions for making trading decisions.
- **Risk Management**: Monitors market conditions and portfolio performance.

## Project Structure

algo_trading/
├── cmd/
│ └── main.go
├── data/
│ └── fetch_data.py
├── internal/
│ ├── app/
│ │ ├── risk/
│ │ │ └── risk_manager.go
│ │ └── trading/
│ │ └── trader.go
│ ├── core/
│ │ ├── entities/
│ │ │ └── data.go
│ │ └── interfaces/
│ │ └── trader.go
│ ├── infrastructure/
│ │ ├── analytics/
│ │ │ └── predictor.go
│ │ ├── db/
│ │ │ └── db.go
│ │ └── di/
│ │ └── container.go
├── model/
│ └── train_model.py
├── service/
│ └── predict_service.py
└── pkg/
└── config/
└── dbconfig.go

kotlin
Copy code

## Setup Instructions

### Prerequisites

- Python 3.8 or higher
- Go 1.16 or higher

### Python Setup

1. **Create a Virtual Environment**:

    ```sh
    python -m venv env
    source env/bin/activate  # On Windows use `env\Scripts\activate`
    ```

2. **Install Dependencies**:

    ```sh
    pip install pandas polygon-api-client scikit-learn joblib flask
    ```

3. **Fetch Historical Market Data**:

    ```sh
    python data/fetch_data.py
    ```

4. **Train the Model**:

    ```sh
    python model/train_model.py
    ```

5. **Start the Prediction Service**:

    ```sh
    python service/predict_service.py
    ```

### Go Setup

1. **Initialize Go Modules**:

    ```sh
    go mod init github.com/yourusername/algo-trading
    ```

2. **Create and Update the Following Files**:

    **`internal/core/entities/data.go`**:
    ```go
    package entities

    type MarketData struct {
        TimeStamp int64
        Symbol    string
        Price     float64
        Volume    float64
        SMA_50    float64
        SMA_200   float64
    }

    type Trade struct {
        ID        string  `json:"ID"`
        Timestamp int64   `json:"Timestamp"`
        Symbol    string  `json:"Symbol"`
        Quantity  float64 `json:"Quantity"`
        Price     float64 `json:"Price"`
        Side      string  `json:"Side"` // "buy" or "sell"
    }

    type Portfolio struct {
        ID         string
        Holdings   map[string]float64 // Symbol -> Quantity
        Cash       float64
        TotalValue float64
    }
    ```

    **`internal/infrastructure/analytics/predictor.go`**:
    ```go
    package analytics

    import (
        "bytes"
        "encoding/json"
        "net/http"
        "errors"
        "github.com/yourusername/algo-trading/internal/core/entities"
    )

    type PredictorService struct {
        APIEndpoint string
    }

    func NewPredictorService(endpoint string) *PredictorService {
        return &PredictorService{APIEndpoint: endpoint}
    }

    func (s *PredictorService) Predict(data entities.MarketData) (int, error) {
        features := map[string]float64{
            "SMA_50": data.SMA_50,
            "SMA_200": data.SMA_200,
        }

        body, err := json.Marshal(features)
        if err != nil {
            return 0, err
        }

        resp, err := http.Post(s.APIEndpoint, "application/json", bytes.NewBuffer(body))
        if err != nil {
            return 0, err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return 0, errors.New("failed to get a valid response from prediction service")
        }

        var result map[string]int
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            return 0, err
        }

        prediction, ok := result["prediction"]
        if !ok {
            return 0, errors.New("prediction key not found in response")
        }

        return prediction, nil
    }
    ```

    **`internal/app/trading/trader.go`**:
    ```go
    package trading

    import (
        "github.com/yourusername/algo-trading/internal/core/entities"
        "github.com/yourusername/algo-trading/internal/core/interfaces"
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
            SMA_50:  trade.Price * 1.05,  // Example placeholder calculation
            SMA_200: trade.Price * 0.95,  // Example placeholder calculation
        })
        if err != nil {
            return err
        }
        
        // Logic based on prediction
        if prediction == 1 {
            // Execute trade if prediction is positive
            return s.storage.SaveTrade(trade)
        }

        return nil  // Do not execute trade if prediction is not positive
    }

    func (s *TraderService) GetPortfolio() (entities.Portfolio, error) {
        return s.storage.GetPortfolio()
    }
    ```

    **`internal/infrastructure/di/container.go`**:
    ```go
    package di

    import (
        "github.com/yourusername/algo-trading/internal/app/risk"
        "github.com/yourusername/algo-trading/internal/app/trading"
        "github.com/yourusername/algo-trading/internal/infrastructure/db"
        "github.com/yourusername/algo-trading/internal/infrastructure/analytics"
        "github.com/yourusername/algo-trading/pkg/config"
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
    ```

    **`cmd/main.go`**:
    ```go
    package main

    import (
        "encoding/json"
        "log"
        "net/http"
        "github.com/yourusername/algo-trading/internal/core/entities"
        "github.com/yourusername/algo-trading/internal/infrastructure/di"
        "github.com/rs/cors"
    )

    func main() {
        container, err := di.NewContainer()
        if err != nil {
            log.Fatal(err)
        }

        http.HandleFunc("/trade", func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodPost {
                http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
                return
            }

            var trade entities.Trade
            if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
            }

            if err := container.TraderService.ExecuteTrade(trade); err != nil {
                http.Error(w, "Failed to execute trade", http.StatusInternalServerError)
                return
            }

            w.WriteHeader(http.StatusOK)
        })

        http.HandleFunc("/portfolio", func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodGet {
                http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
                return
            }

            portfolio, err := container.TraderService.GetPortfolio()
            if err != nil {
                http.Error(w, "Failed to get portfolio", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            if err := json.NewEncoder(w).Encode(portfolio); err != nil {
                http.Error(w, "Failed to encode portfolio", http.StatusInternalServerError)
                return
            }
        })

        corsHandler := cors.New(cors.Options{
            AllowedOrigins:   []string{"http://localhost:3000"},
            AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
            AllowedHeaders:   []string{"Content-Type"},
            AllowCredentials: true,
        })

        handler := corsHandler.Handler(http.DefaultServeMux)

        log.Println("Starting server on :8080")
        if err := http.ListenAndServe(":8080", handler); err != nil {
            log.Fatal(err)
        }
    }
    ```

### Running the Project

1. **Fetch Data**:

    ```sh
    python data/fetch_data.py
    ```

2. **Train Model**:

    ```sh
    python model/train_model.py
    ```

3. **Start the Prediction Service**:

    ```sh
    python service/predict_service.py
    ```

4. **Start the Go Backend**:

    ```sh
    go run cmd/main.go
    ```

5. **Start the React Development Server**:

    ```sh
    npm start
    ```

## Conclusion

This project sets up an algorithmic trading system with data ingestion, model training, prediction service, and a Go-based trading logic. It uses a virtual environment for Python dependencies and integrates seamlessly with a Go backend. Follow the setup instructions to get the project running locally.
