package entities

type MarketData struct {
	TimeStamp int64
	Symbol    string
	Price     float64
	Volume    float64
	SMA_50    float64 // Add field for 50-day Simple Moving Average
	SMA_200   float64 // Add field for 200-day Simple Moving Average
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
