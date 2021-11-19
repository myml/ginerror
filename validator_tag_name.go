package ginerror

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterTagName 注册一个tag name func到gin的验证引擎
// 用于自动从tag(优先级：uri,form,json)生成验证错误的字段名
// 使用tag做为错误字段名，有利于前端自动化提示错误字段
// 只支持新版本gin使用的 validator/v10
func RegisterTagName() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			for _, tag := range []string{"header", "uri", "form", "json"} {
				v := field.Tag.Get(tag)
				if len(v) > 0 {
					return strings.Split(v, ",")[0]
				}
			}
			return field.Name
		})
	} else {
		panic("Only supports validator v10")
	}
}
