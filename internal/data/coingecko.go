package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ni5arga/stock-tui/internal/models"
)

const coingeckoBase = "https://api.coingecko.com/api/v3"

type CoinGecko struct {
	idMap map[string]string
}

func NewCoinGecko() *CoinGecko {
	return &CoinGecko{
		idMap: map[string]string{
			"BTC":     "bitcoin",
			"BTC-USD": "bitcoin",
			"ETH":     "ethereum",
			"ETH-USD": "ethereum",
			"SOL":     "solana",
			"SOL-USD": "solana",
			"XRP":     "ripple",
			"XRP-USD": "ripple",
			"DOGE":    "dogecoin",
			"ADA":     "cardano",
			"DOT":     "polkadot",
			"AVAX":    "avalanche-2",
			"MATIC":   "matic-network",
			"LINK":    "chainlink",
		},
	}
}

func (c *CoinGecko) Name() string { return "CoinGecko" }

func (c *CoinGecko) symbolToID(symbol string) string {
	sym := strings.ToUpper(strings.TrimSuffix(symbol, "-USD"))
	if id, ok := c.idMap[sym]; ok {
		return id
	}
	return strings.ToLower(sym)
}

func (c *CoinGecko) GetQuotes(symbols []string) ([]models.Quote, error) {
	ids := make([]string, 0, len(symbols))
	symToID := make(map[string]string)
	for _, s := range symbols {
		id := c.symbolToID(s)
		ids = append(ids, id)
		symToID[s] = id
	}

	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd&include_24hr_change=true",
		coingeckoBase, strings.Join(ids, ","))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	body, err := fetch(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	var data map[string]struct {
		USD       float64 `json:"usd"`
		Change24h float64 `json:"usd_24h_change"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	now := time.Now()
	quotes := make([]models.Quote, 0, len(symbols))
	for _, sym := range symbols {
		id := symToID[sym]
		if d, ok := data[id]; ok {
			quotes = append(quotes, models.Quote{
				Symbol:      sym,
				Price:       d.USD,
				ChangePct:   d.Change24h,
				LastUpdated: now,
			})
		}
	}

	return quotes, nil
}

func (c *CoinGecko) GetHistory(symbol string, tr models.TimeRange) ([]models.Candle, error) {
	id := c.symbolToID(symbol)

	var days string
	switch tr {
	case models.Range1H:
		days = "1"
	case models.Range24H:
		days = "1"
	case models.Range7D:
		days = "7"
	case models.Range30D:
		days = "30"
	default:
		days = "1"
	}

	url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=usd&days=%s", coingeckoBase, id, days)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	body, err := fetch(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	var data struct {
		Prices [][]float64 `json:"prices"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	candles := make([]models.Candle, 0, len(data.Prices))
	for _, p := range data.Prices {
		if len(p) < 2 {
			continue
		}
		ts := time.UnixMilli(int64(p[0]))
		price := p[1]
		candles = append(candles, models.Candle{
			Timestamp: ts,
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
		})
	}

	return candles, nil
}
