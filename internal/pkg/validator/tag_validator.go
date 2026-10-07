package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
	entranslation "github.com/go-playground/validator/v10/translations/en"
	rutranslation "github.com/go-playground/validator/v10/translations/ru"
)

// Ensure tag validator implements interface.
var _ Validator = (*TagValidator)(nil)

// TagValidator implementats Validator interface.
// It validates structs by tags.
type TagValidator struct {
	validatorInstance *govalidator.Validate
	translator        ut.Translator
}

// NewTagValidator returns new instance of TagValidator.
func NewTagValidator() (*TagValidator, error) {
	enTranslator := en.New()
	uni := ut.New(enTranslator, enTranslator)
	trans, _ := uni.GetTranslator("en")

	validate := govalidator.New(govalidator.WithRequiredStructEnabled())
	err := entranslation.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, fmt.Errorf("register translation: %w", err)
	}

	validInst := &TagValidator{validate, trans}
	return validInst, nil
}

// NewRuTagValidator returns new instance of TagValidator with ru error translation.
func NewRuTagValidator() (*TagValidator, error) {
	ruTranslator := ru.New()
	uni := ut.New(ruTranslator, ruTranslator)
	trans, _ := uni.GetTranslator("ru")

	validate := govalidator.New(govalidator.WithRequiredStructEnabled())
	err := rutranslation.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, fmt.Errorf("register translation: %w", err)
	}

	validInst := &TagValidator{validate, trans}
	return validInst, nil
}

// Validate validates given struct s (s is a pointer to struct).
func (v TagValidator) Validate(s any) error {
	err := v.validatorInstance.Struct(s)
	if err == nil { // NOT err
		return nil
	}

	var validateErrors govalidator.ValidationErrors
	if !errors.As(err, &validateErrors) {
		return err
	}
	// handle error messages
	rawTranstaledMap := validateErrors.Translate(v.translator)

	// sort out errors and concat them into string
	transtaledStringSlice := make([]string, 0, len(rawTranstaledMap))
	for _, v := range rawTranstaledMap {
		transtaledStringSlice = append(transtaledStringSlice, strings.ToLower(v))
	}

	return fmt.Errorf("%w: %s", ErrValidate, strings.Join(transtaledStringSlice, "\n"))
}
