package validation

import (
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
}
