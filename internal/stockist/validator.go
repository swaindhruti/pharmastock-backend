package stockist

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateStockist(stockist *Stockist) error {
	return validate.Struct(stockist)
}
