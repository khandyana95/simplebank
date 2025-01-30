package util

const (
	USD = "USD"
	INR = "INR"
	CAD = "CAD"
)

func ValidateCurrency(currency string) bool {
	switch currency {
	case USD, INR, CAD:
		return true
	}

	return false
}
