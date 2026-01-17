package data

import (
	"math"
	"math/rand"
	"time"

	"github.com/ni5arga/stock-tui/internal/models"
)

type Simulator struct {
	basePrices map[string]float64
}

func NewSimulator() *Simulator {
	return &Simulator{
		basePrices: map[string]float64{
			"BTC-USD": 95000.0,
			"ETH-USD": 3400.0,
			"AAPL":    225.0,
			"GOOGL":   175.0,
			"TSLA":    240.0,
		},
	}
}

func (s *Simulator) Name() string { return "Simulator" }

func (s *Simulator) GetQuotes(symbols []string) ([]models.Quote, error) {
	var quotes []models.Quote
	now := time.Now()

	for _, sym := range symbols {
		base, ok := s.basePrices[sym]
		if !ok {
			base = 100.0 // Default for unknown symbols
		}

		// Random walk
		volatility := base * 0.02
		change := (rand.Float64() - 0.5) * volatility
		current := base + change
		pct := (change / base) * 100

		quotes = append(quotes, models.Quote{
			Symbol:      sym,
			Price:       current,
			ChangePct:   pct,
			LastUpdated: now,
		})
	}
	return quotes, nil
}

func (s *Simulator) GetHistory(symbol string, tr models.TimeRange) ([]models.Candle, error) {
	var points int
	var duration time.Duration

	switch tr {
	case models.Range1H:
		points = 60
		duration = time.Minute
	case models.Range7D:
		points = 84 // 4-hour intervals roughly
		duration = 4 * time.Hour
	case models.Range30D:
		points = 30
		duration = 24 * time.Hour
	default: // 24H
		points = 48 // 30-min intervals
		duration = 30 * time.Minute
	}

	base, ok := s.basePrices[symbol]
	if !ok {
		base = 100.0
	}

	candles := make([]models.Candle, points)
	currentPrice := base

	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(points) * duration)

	for i := 0; i < points; i++ {
		// Generate random candle data
		volatility := currentPrice * 0.01
		open := currentPrice
		close := currentPrice + (rand.Float64()-0.5)*volatility

		high := math.Max(open, close) + rand.Float64()*volatility*0.5
		low := math.Min(open, close) - rand.Float64()*volatility*0.5

		candles[i] = models.Candle{
			Timestamp: startTime.Add(time.Duration(i) * duration),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    rand.Float64() * 1000,
		}

		currentPrice = close
	}

	return candles, nil
}
