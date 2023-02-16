package validatorMng

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
)

type ValidatorMng struct{}

var validate = validator.New()

func GetError(s interface{}) error {
	err := validate.Struct(s)
	validationErrors := err.(validator.ValidationErrors)
	log.Println("validationErrors", validationErrors)

	for k := range validationErrors {
		log.Println("validationErrors", k, validationErrors[k])
	}

	if err == nil {
		return nil
	}
	Report(err)
	return err
}

func Report(errs error) {
	for _, err := range errs.(validator.ValidationErrors) {
		fmt.Println("---------------------------------------")
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
		fmt.Println("---------------------------------------")
	}
}
