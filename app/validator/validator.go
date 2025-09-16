package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate);
	if ok {
		_ = v.RegisterValidation("pwdmin", pwdminFunc)
	}
}

func pwdminFunc(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return len(val) >= 8 || val == ""
}