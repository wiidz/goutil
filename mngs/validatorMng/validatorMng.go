package validatorMng

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

func GetError(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	Report(err)
	return err
}

func Report(errs error) {
	for _, err := range errs.(validator.ValidationErrors) {
		fmt.Println()
		fmt.Println(err.Namespace())
		fmt.Println(err.Field())
		fmt.Println(err.StructNamespace())
		fmt.Println(err.StructField())
		fmt.Println(err.Tag())
		fmt.Println(err.ActualTag())
		fmt.Println(err.Kind())
		fmt.Println(err.Type())
		fmt.Println(err.Value())
		fmt.Println(err.Param())
		fmt.Println()
	}
}
