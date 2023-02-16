package validatorMng

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"log"
	"reflect"
	"strings"
)

type ValidatorMng struct{}

var validate = validator.New()
var trans ut.Translator

func init() {
	en := en.New() //英文翻译器
	zh := zh.New() //中文翻译器

	// 第一个参数是必填，如果没有其他的语言设置，就用这第一个
	// 后面的参数是支持多语言环境（
	// uni := ut.New(en, en) 也是可以的
	// uni := ut.New(en, zh, tw)
	uni := ut.New(en, zh)
	trans, _ = uni.GetTranslator("zh") //获取需要的语言

	zhtrans.RegisterDefaultTranslations(validate, trans)
}

func GetError(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	return TranslateOne(s, err)
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

// GetTranslate 获取带翻译器的validate
//func GetTranslate(){
//
//}

// TranslateOne 翻译一下
//func TranslateOne(errs error) (err error) {
//
//	//【1】提取错误
//	validationErrors := errs.(validator.ValidationErrors)
//	if len(validationErrors) == 0 {
//		return errs
//	}
//
//	//【2】只取一号错误
//	return NewValidatorErr(validationErrors[0])
//}

// TranslateOne 翻译一下
func TranslateOne(params interface{}, errs error) (err error) {

	translatedErrs := errs.(validator.ValidationErrors).Translate(trans)

	// 只取一号错误
	for k := range translatedErrs {
		log.Println("k", k)
		//【2】获取字段定义
		structType := reflect.TypeOf(params)
		log.Println("structType", structType)

		tempArr := typeHelper.ExplodeStr(k, ".") // 这个时候还是MyStruct.TrueName,所以要提取TrueName
		if len(tempArr) < 2 {
			return errs
		}
		field, _ := structType.Elem().FieldByName(tempArr[1])

		log.Println("field", field)
		cnTag := field.Tag.Get("cn") // 如果定义了中文
		log.Println("cnTag", cnTag)

		//【3】替换字段名
		errStr := translatedErrs[k]
		if cnTag != "" {
			errStr = strings.Replace(errStr, k, cnTag, 1)
		}

		return errors.New(errStr)
	}

	return
}
