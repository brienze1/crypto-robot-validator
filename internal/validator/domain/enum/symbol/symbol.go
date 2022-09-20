package symbol

type Symbol string

const (
	Bitcoin    Symbol = "BTC"
	Brl        Symbol = "BRL"
	BitcoinBRL Symbol = "BTCBRL"
)

func (s Symbol) Name() string {
	return string(s)
}
