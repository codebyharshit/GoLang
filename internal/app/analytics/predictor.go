package analytics

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/codebyharshit/real-time-analytics/internal/core/entities"
)

type PredictorService struct {
	APIEndpoint string
}

func NewPredictorService(endpoint string) *PredictorService {
	return &PredictorService{APIEndpoint: endpoint}
}

func (s *PredictorService) Predict(data entities.MarketData) (int, error) {
	features := map[string]float64{
		"SMA_50":  data.SMA_50,
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
