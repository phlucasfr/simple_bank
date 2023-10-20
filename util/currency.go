package util

const (
	USD = "USD"
	BRL = "BRL"
	EUR = "EUR"
)

// IsSupportedCurrency returns true or false
func IsSupportedCurrency(currency string) bool {

	switch currency {
	case USD, BRL, EUR:
		return true
	}

	return false
}
