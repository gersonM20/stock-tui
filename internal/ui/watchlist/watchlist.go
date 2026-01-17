package watchlist

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ni5arga/stock-tui/internal/models"
	"github.com/ni5arga/stock-tui/internal/ui/styles"
)

type Model struct {
	list   list.Model
	width  int
	height int
}

type item struct {
	symbol    string
	price     float64
	changePct float64
}

func (i item) Title() string       { return i.symbol }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.symbol }

func New(symbols []string) Model {
	items := make([]list.Item, len(symbols))
	for i, s := range symbols {
		items[i] = item{symbol: s}
	}

	l := list.New(items, newDelegate(), 0, 0)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(true)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()

	return Model{list: l}
}

type delegate struct{}

func newDelegate() delegate { return delegate{} }

func (d delegate) Height() int                               { return 1 }
func (d delegate) Spacing() int                              { return 0 }
func (d delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(item)
	if !ok {
		return
	}

	// Dynamic widths based on list width
	totalW := m.Width()
	symW := 14
	priceW := 12
	pctW := 9

	if totalW > 40 {
		symW = min(20, totalW-priceW-pctW-2)
	}

	// Symbol - truncate if needed
	sym := it.symbol
	if len(sym) > symW {
		sym = sym[:symW-1] + "…"
	}
	symStr := fmt.Sprintf("%-*s", symW, sym)

	// Price
	var priceStr string
	if it.price == 0 {
		priceStr = fmt.Sprintf("%*s", priceW, "—")
	} else if it.price >= 1000 {
		priceStr = fmt.Sprintf("%*.0f", priceW, it.price)
	} else {
		priceStr = fmt.Sprintf("%*.2f", priceW, it.price)
	}

	// Percent change
	var pctStr string
	if it.price == 0 {
		pctStr = fmt.Sprintf("%*s", pctW, "—")
	} else {
		pctStr = fmt.Sprintf("%+*.2f%%", pctW-1, it.changePct)
	}

	// Style based on selection and trend
	selected := index == m.Index()

	if selected {
		row := fmt.Sprintf("%s %s %s", symStr, priceStr, pctStr)
		fmt.Fprint(w, styles.SelectedItem.Render(row))
	} else {
		symStyled := lipgloss.NewStyle().Foreground(styles.ColorText).Render(symStr)
		priceStyled := lipgloss.NewStyle().Foreground(styles.ColorText).Render(priceStr)

		pctStyle := styles.PositiveChange
		if it.changePct < 0 {
			pctStyle = styles.NegativeChange
		}
		pctStyled := pctStyle.Render(pctStr)

		fmt.Fprint(w, fmt.Sprintf(" %s %s %s", symStyled, priceStyled, pctStyled))
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return styles.Pane.
		Width(m.width).
		Height(m.height).
		Render(m.list.View())
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.list.SetSize(w-4, h-2)
}

func (m *Model) UpdateQuotes(quotes []models.Quote) {
	items := m.list.Items()
	updated := make([]list.Item, len(items))

	qmap := make(map[string]models.Quote, len(quotes))
	for _, q := range quotes {
		qmap[q.Symbol] = q
	}

	for i, it := range items {
		curr := it.(item)
		if q, ok := qmap[curr.symbol]; ok {
			curr.price = q.Price
			curr.changePct = q.ChangePct
		}
		updated[i] = curr
	}

	m.list.SetItems(updated)
}

// UpdatePriceChange updates change % for a symbol based on historical data
func (m *Model) UpdatePriceChange(symbol string, currentPrice, startPrice float64) {
	items := m.list.Items()
	for i, it := range items {
		curr := it.(item)
		if curr.symbol == symbol {
			curr.price = currentPrice
			if startPrice > 0 {
				curr.changePct = ((currentPrice - startPrice) / startPrice) * 100
			}
			items[i] = curr
			break
		}
	}
	m.list.SetItems(items)
}

func (m Model) SelectedSymbol() string {
	if it, ok := m.list.SelectedItem().(item); ok {
		return it.symbol
	}
	return ""
}
