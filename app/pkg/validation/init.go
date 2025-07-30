package validation

import (
	"slices"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

var Validate *validator.Validate

func Init() {
	Validate = validator.New()
	Validate.RegisterValidation("YYYY-MM-DD_dateFormat", func(fl validator.FieldLevel) bool {
		_, err := time.Parse("2006-01-02", fl.Field().String())
		return err == nil
	})
	Validate.RegisterValidation("size", func(fl validator.FieldLevel) bool {
		value := strings.ToLower(fl.Field().String())
		allowedValues := []string{"small", "medium", "large"}
		return slices.Contains(allowedValues, value)
	})
	Validate.RegisterValidation("flavor", func(fl validator.FieldLevel) bool {
		value := strings.ToLower(fl.Field().String())
		allowedValues := []string{"jagung bakar", "rumput laut", "original", "jagung manis", "keju asin", "keju manis", "pedas"}
		return slices.Contains(allowedValues, value)
	})
	Validate.RegisterValidation("type", func(fl validator.FieldLevel) bool {
		value := strings.ToLower(fl.Field().String())
		return value == "keripik pangsit"
	})
	Validate.RegisterValidation("price", func(fl validator.FieldLevel) bool {
		price := fl.Field().Float()
		switch fl.Parent().FieldByName("Size").String() {
		case "small":
			return price == 10000
		case "medium":
			return price == 25000
		case "large":
			return price == 35000
		default:
			return false
		}
	})
}
