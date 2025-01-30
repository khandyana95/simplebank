package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/khandyan95/simplebank/util"
)

var currencyValidator validator.Func = func(f validator.FieldLevel) bool {
	if currency, ok := f.Field().Interface().(string); ok {
		return util.ValidateCurrency(currency)
	}

	return false
}
