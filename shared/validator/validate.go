package validator

import (
	dto "shared/models"

	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tr_translations "github.com/go-playground/validator/v10/translations/tr"
)

var (
	Validate *validator.Validate
	Trans    ut.Translator
)

func Init() {
	turkish := tr.New()
	uni := ut.New(turkish, turkish)
	Trans, _ = uni.GetTranslator("tr")

	Validate = validator.New()
	tr_translations.RegisterDefaultTranslations(Validate, Trans)
}

func CheckStruct(s interface{}) *dto.GeneralResponse[any] {
	err := Validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = e.Translate(Trans)
	}

	return &dto.GeneralResponse[any]{
		Status:  false,
		Message: "Validasyon hatası oluştu.",
		Errors:  errors,
		Code:    400,
	}
}
