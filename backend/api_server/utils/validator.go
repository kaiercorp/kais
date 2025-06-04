package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"api_server/logger"
)

type Validator struct{}

var validatorInstance *Validator

func NewValidator() *Validator {
	if validatorInstance == nil {
		validatorInstance = &Validator{}
	}

	return validatorInstance
}

func (vldt *Validator) RegisterValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("sampleValidation", sampleValidation)
		if err != nil {
			logger.Error("Register Validation Failed", err)
		}
	}
}

var sampleValidation validator.Func = func(fl validator.FieldLevel) bool {
	return true
}
