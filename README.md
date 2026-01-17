# stock-tui

Real-time stock and cryptocurrency tracker for the terminal.

![screenshot](screenshots/stock-tui.png)

## Features

- Real-time price tracking for stocks and cryptocurrencies
- Multiple data providers (CoinGecko, Yahoo Finance, or combined)
- Historical price charts with multiple time ranges
- Sparkline visualization
- Keyboard-driven interface with Vim-style navigation

## Installation

```bash
go install github.com/ni5arga/stock-tui/cmd/stock-tui@latest
```

Or build from source:

```bash
git clone https://github.com/ni5arga/stock-tui.git
cd stock-tui
go build ./cmd/stock-tui
./stock-tui
```

## Configuration

The app looks for config in this order:
1. `~/.config/stock-tui/config.toml` — user config
2. `./config.toml` — current directory

A sample `config.toml` is included in the repo. To use it system-wide:

```bash
mkdir -p ~/.config/stock-tui
cp config.toml ~/.config/stock-tui/
```

**Example config.toml:**

```toml
# Data provider: "simulator", "coingecko", "yahoo", or "multi" (default)
provider = "multi"

# Refresh interval
refresh_interval = "5s"

# Default chart range: "1H", "24H", "7D", "30D"
default_range = "24H"

# Watchlist symbols
# Crypto: use -USD suffix (BTC-USD, ETH-USD)
# Stocks: use ticker (AAPL, GOOGL)
symbols = [
    "BTC-USD",
    "ETH-USD",
    "SOL-USD",
    "AAPL",
    "GOOGL",
    "TSLA",
    "MSFT",
    "NVDA"
]
```

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down in watchlist |
| `k` / `↑` | Move up in watchlist |
| `Tab` | Cycle time range |
| `1` | 1 hour range |
| `2` | 24 hour range |
| `3` | 7 day range |
| `4` | 30 day range |
| `c` | Cycle chart type (Line/Area/Bars/Dots) |
| `r` | Refresh data |
| `?` | Toggle help |
| `q` | Quit |

## Data Providers

| Provider | Assets | API Key |
|----------|--------|---------|
| `simulator` | Fake data | None |
| `coingecko` | Crypto | None (free tier) |
| `yahoo` | Stocks | None (unofficial) |
| `multi` | Both | None |

> **Note**: Yahoo Finance API is unofficial and may have rate limits.
> CoinGecko free tier allows ~10-30 requests/minute.

## Supported Platforms

- Linux
- macOS
- Windows

## Architecture

```
cmd/stock-tui/       Entry point
internal/
├── app/             Bubble Tea model
├── config/          Viper configuration
├── data/            Provider implementations
├── models/          Domain types
└── ui/
    ├── chart/       Price chart component
    ├── footer/      Status bar
    ├── help/        Help overlay
    ├── modal/       Generic modal
    ├── styles/      Lip Gloss styles
    └── watchlist/   Symbol list
```

## Development

```bash
# Run
go run ./cmd/stock-tui

# Build
go build ./cmd/stock-tui

# Test
go test ./...

# Lint
go vet ./...
```

## License

MIT
