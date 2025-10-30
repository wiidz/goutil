package paramHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/validatorMng"
	"github.com/wiidz/goutil/structs/networkStruct"
)

type (
	// BuildParamsOption allows callers to customize the build pipeline.
	BuildParamsOption func(*buildParamsOptions)
	// ParamValidator is a validation step executed after the params are filled.
	ParamValidator func(networkStruct.ParamsInterface) error
	// ParamMutator is a processing step executed after validation succeeds.
	ParamMutator func(networkStruct.ParamsInterface) error
)

type buildParamsOptions struct {
	skipValidation bool
	skipHandle     bool
	validators     []ParamValidator
	mutators       []ParamMutator
}

var defaultParamValidator ParamValidator = func(params networkStruct.ParamsInterface) error {
	return validatorMng.GetError(params)
}

func newBuildParamsOptions() buildParamsOptions {
	return buildParamsOptions{
		validators: []ParamValidator{defaultParamValidator},
		mutators:   []ParamMutator{handleParams},
	}
}

// WithSkipValidation skips all validation steps (including defaults).
func WithSkipValidation() BuildParamsOption {
	return func(o *buildParamsOptions) {
		o.skipValidation = true
	}
}

// WithValidators appends additional validation steps. They run after the default validator.
func WithValidators(validators ...ParamValidator) BuildParamsOption {
	return func(o *buildParamsOptions) {
		for _, validator := range validators {
			if validator != nil {
				o.validators = append(o.validators, validator)
			}
		}
	}
}

// WithSkipHandle skips all post-validation handlers (including defaults).
func WithSkipHandle() BuildParamsOption {
	return func(o *buildParamsOptions) {
		o.skipHandle = true
	}
}

// WithMutators appends additional post-validation handlers. They run after the default handler.
func WithMutators(mutators ...ParamMutator) BuildParamsOption {
	return func(o *buildParamsOptions) {
		for _, mutator := range mutators {
			if mutator != nil {
				o.mutators = append(o.mutators, mutator)
			}
		}
	}
}

// BuildParams 构建参数
func BuildParams(r *http.Request, params networkStruct.ParamsInterface, contentType networkStruct.ContentType, opts ...BuildParamsOption) error {
	if r == nil {
		return errors.New("build params: request is nil")
	}
	if params == nil {
		return errors.New("build params: params is nil")
	}

	cfg := newBuildParamsOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if err := fillParams(r, params, contentType); err != nil {
		return fmt.Errorf("build params: fill: %w", err)
	}

	if !cfg.skipValidation {
		for _, validator := range cfg.validators {
			if validator == nil {
				continue
			}
			if err := validator(params); err != nil {
				return fmt.Errorf("build params: validate: %w", err)
			}
		}
	}

	if cfg.skipHandle {
		return nil
	}

	for _, mutator := range cfg.mutators {
		if mutator == nil {
			continue
		}
		if err := mutator(params); err != nil {
			return fmt.Errorf("build params: mutate: %w", err)
		}
	}

	return nil
}

// fillParams 把请求中的参数填充到结构体
func fillParams(r *http.Request, params networkStruct.ParamsInterface, contentType networkStruct.ContentType) (err error) {
	switch contentType {
	case networkStruct.Query:
		if err = r.ParseForm(); err != nil {
			return
		}
		queryMap := map[string]interface{}{}
		for k, v := range r.URL.Query() {
			if len(v) == 1 {
				queryMap[k] = v[0]
			} else {
				queryMap[k] = v
			}
		}
		params.SetRawMap(queryMap)
		b, _ := json.Marshal(queryMap)
		err = json.Unmarshal(b, params)

	case networkStruct.BodyJson:
		buf, _ := ioutil.ReadAll(r.Body)
		jsonStr := string(buf)
		jsonMap := typeHelper.JsonDecodeMap(jsonStr)
		params.SetRawMap(jsonMap)
		err = typeHelper.JsonDecodeWithStruct(jsonStr, params)
		if err != nil {
			log.Println("err", err, params)
			return
		}

	case networkStruct.BodyForm:
		formMap := make(map[string]interface{})
		if err = r.ParseForm(); err != nil {
			return
		}
		for k, v := range r.PostForm {
			if len(v) == 1 {
				formMap[k] = v[0]
			} else {
				formMap[k] = v
			}
		}
		params.SetRawMap(formMap)

	case networkStruct.XWWWForm:
		if err = r.ParseForm(); err != nil {
			return
		}
		formMap := make(map[string]interface{})
		for k, v := range r.PostForm {
			if len(v) == 1 {
				formMap[k] = v[0]
			} else {
				formMap[k] = v
			}
		}
		params.SetRawMap(formMap)
		b, _ := json.Marshal(formMap)
		err = json.Unmarshal(b, params)

	default:
		err = errors.New("未能匹配数据类型")
	}
	return
}

// handleParams 处理和填充数据
func handleParams(params networkStruct.ParamsInterface) (err error) {
	structType := reflect.TypeOf(params)
	structValues := reflect.ValueOf(params)
	rawMap := params.GetRawMap()

	condition := map[string]interface{}{}
	value := map[string]interface{}{}
	etc := map[string]interface{}{}

	for i := 0; i < structType.Elem().NumField(); i++ {
		field := structType.Elem().Field(i)
		fieldType := field.Type

		jsonTag := field.Tag.Get("json")
		urlTag := field.Tag.Get("url")
		fieldName := field.Tag.Get("field")

		belong := field.Tag.Get("belong")
		kind := field.Tag.Get("kind")

		defaultValue := field.Tag.Get("default")

		if jsonTag == "" && urlTag == "" {
			continue
		}

		if fieldName == "" {
			if jsonTag != "" {
				fieldName = jsonTag
			} else if urlTag != "" {
				fieldName = urlTag
			}
		}

		currentValue := reflect.Indirect(structValues).FieldByName(field.Name)
		if fieldType.Kind() == reflect.Struct {
		} else if fieldType.Kind() == reflect.Slice {
		} else if fieldType.Kind() == reflect.Map {
		} else if currentValue.Interface() == reflect.Zero(fieldType).Interface() {
			if defaultValue != "" {
				var formattedDefaultValue interface{}
				formattedDefaultValue, err = getFormattedValue(field.Type.String(), defaultValue)
				if err != nil {
					return
				}
				currentValue = reflect.ValueOf(formattedDefaultValue)
				structValues.Elem().Field(i).Set(currentValue)
			} else if _, ok := rawMap[fieldName]; !ok {
				continue
			}
		}

		var formattedValue interface{}
		formattedValue, err = getFormattedValue(field.Type.String(), currentValue.Interface())
		if err != nil {
			log.Println("field.Type.String()", field.Type.String(), err)
			return
		}

		switch belong {
		case "condition":
			switch kind {
			case "like":
				condition[fieldName] = []interface{}{"like", "%" + formattedValue.(string) + "%"}
			case "between":
				tempSlice := typeHelper.Explode(formattedValue.(string), ",")
				if len(tempSlice) == 2 {
					condition[fieldName] = []interface{}{"between", tempSlice[0], tempSlice[1]}
				}
			case "in":
				condition[fieldName] = []interface{}{"in"}
				condition[fieldName] = append(condition[fieldName].([]interface{}), formattedValue)
			case "not in":
				condition[fieldName] = []interface{}{"not in"}
				condition[fieldName] = append(condition[fieldName].([]interface{}), formattedValue)
			case "!=":
				condition[fieldName] = []interface{}{"!="}
				condition[fieldName] = append(condition[fieldName].([]interface{}), formattedValue)
			default:
				condition[fieldName] = formattedValue
			}
		case "value":
			value[fieldName] = formattedValue
		case "etc":
			etc[fieldName] = formattedValue
		}
	}

	params.SetCondition(condition)
	params.SetValue(value)
	params.SetEtc(etc)

	if pageNow, ok := etc["page_now"].(int); ok {
		params.SetPageNow(pageNow)
	} else {
		params.SetPageNow(0)
	}
	if pageSize, ok := etc["page_size"].(int); ok {
		params.SetPageSize(pageSize)
	} else {
		params.SetPageSize(10)
	}
	if order, ok := etc["order"].(string); ok {
		params.SetOrder(order)
	} else {
		params.SetOrder("id asc")
	}
	return
}

func getFormattedValue(t string, value interface{}) (data interface{}, err error) {
	switch t {
	case "string":
		data = typeHelper.ForceString(value)
	case "int":
		data = typeHelper.ForceInt(value)
	case "int8":
		data = typeHelper.ForceInt8(value)
	case "uint64":
		data = typeHelper.ForceUint64(value)
	case "float64":
		data = typeHelper.ForceFloat64(value)
	case "[]int":
		data = typeHelper.ForceIntSlice(value)
	case "[]int8":
		data = typeHelper.ForceInt8Slice(value)
	case "[]uint64":
		data = typeHelper.ForceUint64Slice(value)
	case "[]float64":
		data = typeHelper.ForceFloat64Slice(value)
	case "[]string":
		data = typeHelper.ForceStrSlice(value)
	default:
		data = value
	}
	return
}
