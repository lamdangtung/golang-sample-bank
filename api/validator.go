package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/lamdangtung/golang-sample-bank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		if util.IsSupportedCurrency(currency) {
			return true
		}
	}
	return false
}
