package validator

import (
	"regexp"
	"strings"

	"gopkg.in/go-playground/validator.v9"

	"ImaginatoGolangTestTask/shared/common"
)

type IValidatorService interface {
	ValidateStruct(req interface{}, name string) (string, bool)
}
type Validator struct{}

func NewValidatorService() IValidatorService {
	return &Validator{}
}

func (av *Validator) ValidateStruct(admin interface{}, key string) (string, bool) {
	validate := validator.New()
	err := validate.Struct(admin)
	var errorString string
	if err != nil {
		valErrs := err.(validator.ValidationErrors)
		for _, v := range valErrs {
			fieldName := strings.Replace(strings.Replace(v.Namespace(), key+".", "", 1), ".", " ", 3)
			reg, _ := regexp.Compile("[^A-Z`[]]+")
			fieldName = strings.Replace(reg.ReplaceAllString(fieldName, ""), "[", "", 2)
			errorString = common.GetError(fieldName, v.Tag())
			return errorString, false
		}
	}
	return errorString, true
}
