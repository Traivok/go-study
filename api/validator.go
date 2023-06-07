package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/traivok/go-study/util"
)

var validateCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}

	return false
}
