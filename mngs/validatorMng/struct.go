package validatorMng

import (
	"github.com/go-playground/validator/v10"
)

type ValidatorErr struct {
	OriginalErr   validator.FieldError
	EnField       string `json:"en_field"`       // 英文字段
	CnField       string `json:"cn_field"`       // 中文字段
	ValidationTag string `json:"validation_tag"` // 验证方法
	Format        string `json:"format"`         // 类型验证
}

// NewValidatorErr 构建一个验证器错误
func NewValidatorErr(originErr validator.FieldError) *ValidatorErr {
	return &ValidatorErr{
		OriginalErr:   originErr,
		EnField:       originErr.Field(),
		CnField:       originErr.StructField(), // 记得定义一下
		ValidationTag: originErr.Tag(),
		Format:        originErr.ActualTag(),
	}
}

func (obj *ValidatorErr) Error() (errStr string) {

	//【1】确定字段名称
	fieldName := obj.CnField
	if fieldName == "" {
		fieldName = obj.EnField
	}

	//【2】判断错误类型
	if obj.ValidationTag == "required" {
		errStr = fieldName + "是必填项"
	} else {
		errStr = fieldName + "参数错误"
	}

	return
}
