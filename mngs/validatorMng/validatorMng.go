package validatorMng

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type ValidatorMng struct{}

var validate = validator.New()

func GetError(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	return TranslateOne(err)
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

// TranslateOne 翻译一下
func TranslateOne(errs error) (err error) {

	//【1】提取错误
	validationErrors := errs.(validator.ValidationErrors)
	if len(validationErrors) == 0 {
		return errs
	}

	//【2】只取一号错误
	return NewValidatorErr(validationErrors[0])
}
