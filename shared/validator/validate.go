package validator

import (
	dto "shared/models"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	tr_translations "github.com/go-playground/validator/v10/translations/tr"
)

var (
	Validate *validator.Validate
	uni      *ut.UniversalTranslator
)

func Init() {
	turkish := tr.New()
	english := en.New()

	// Tüm dilleri içeren bir UniversalTranslator oluştur
	uni = ut.New(english, english, turkish)
	Validate = validator.New()

	// Tüm dillerin çeviri kayıtlarını validator'a önceden yap
	trTrans, _ := uni.GetTranslator("tr")
	tr_translations.RegisterDefaultTranslations(Validate, trTrans)

	enTrans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(Validate, enTrans)
}

func CheckStruct(s interface{}, lang string) *dto.GeneralResponse[any] {
	err := Validate.Struct(s)
	if err == nil {
		return nil
	}

	// İstenen dile ait translator'ı çalışma anında seç
	trans, _ := uni.GetTranslator(lang)

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = e.Translate(trans)
	}

	msg := "Validasyon hatası oluştu."
	if lang == "en" {
		msg = "Validation error occurred."
	}

	return &dto.GeneralResponse[any]{
		Status:  false,
		Message: msg,
		Errors:  errors,
		Code:    400,
	}
}
