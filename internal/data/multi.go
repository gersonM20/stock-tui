package data

import (
	"strings"
	"sync"

	"github.com/ni5arga/stock-tui/internal/models"
)

type Multi struct {
	crypto Provider
	stocks Provider
}

func NewMulti() *Multi {
	return &Multi{
		crypto: NewCoinGecko(),
		stocks: NewYahoo(),
	}
}

func (m *Multi) Name() string { return "Multi (CoinGecko + Yahoo)" }

func (m *Multi) isCrypto(symbol string) bool {
	sym := strings.ToUpper(symbol)
	cryptoSymbols := map[string]bool{
		"BTC": true, "BTC-USD": true,
		"ETH": true, "ETH-USD": true,
		"SOL": true, "SOL-USD": true,
		"XRP": true, "XRP-USD": true,
		"DOGE": true, "DOGE-USD": true,
		"ADA": true, "ADA-USD": true,
		"DOT": true, "DOT-USD": true,
		"AVAX": true, "AVAX-USD": true,
		"MATIC": true, "MATIC-USD": true,
		"LINK": true, "LINK-USD": true,
	}
	return cryptoSymbols[sym]
}

func (m *Multi) GetQuotes(symbols []string) ([]models.Quote, error) {
	var cryptoSyms, stockSyms []string
	for _, s := range symbols {
		if m.isCrypto(s) {
			cryptoSyms = append(cryptoSyms, s)
		} else {
			stockSyms = append(stockSyms, s)
		}
	}

	var wg sync.WaitGroup
	var cryptoQuotes, stockQuotes []models.Quote
	var cryptoErr, stockErr error

	if len(cryptoSyms) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cryptoQuotes, cryptoErr = m.crypto.GetQuotes(cryptoSyms)
		}()
	}

	if len(stockSyms) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stockQuotes, stockErr = m.stocks.GetQuotes(stockSyms)
		}()
	}

	wg.Wait()

	// Return partial results even if one fails
	var quotes []models.Quote
	quotes = append(quotes, cryptoQuotes...)
	quotes = append(quotes, stockQuotes...)

	// Prioritize returning data; only error if both failed
	if len(quotes) == 0 {
		if cryptoErr != nil {
			return nil, cryptoErr
		}
		if stockErr != nil {
			return nil, stockErr
		}
	}

	return quotes, nil
}

func (m *Multi) GetHistory(symbol string, tr models.TimeRange) ([]models.Candle, error) {
	if m.isCrypto(symbol) {
		return m.crypto.GetHistory(symbol, tr)
	}
	return m.stocks.GetHistory(symbol, tr)
}
