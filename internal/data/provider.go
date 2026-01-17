package data

import (
	"fmt"

	"github.com/ni5arga/stock-tui/internal/models"
)

// Provider defines the interface for data sources.
type Provider interface {
	Name() string
	GetQuotes(symbols []string) ([]models.Quote, error)
	GetHistory(symbol string, tr models.TimeRange) ([]models.Candle, error)
}

// NewProvider returns the requested provider implementation.
func NewProvider(name string) (Provider, error) {
	switch name {
	case "simulator":
		return NewSimulator(), nil
	case "coingecko":
		return NewCoinGecko(), nil
	case "yahoo":
		return NewYahoo(), nil
	case "multi", "auto":
		return NewMulti(), nil
	default:
		return NewMulti(), fmt.Errorf("unknown provider %q, using multi", name)
	}
}
