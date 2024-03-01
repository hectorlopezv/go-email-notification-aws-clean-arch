package config

import "github.com/go-playground/validator/v10"



type (
    ScanAvRequest struct {
      Paths []string `json:"paths" validate:"required"`
    }

    ErrorResponse struct {
        Error       bool
        FailedField string
        Tag         string
        Value       any
    }

    XValidator struct {
        Validator *validator.Validate
    }

    GlobalErrorHandlerResp struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }
)
var validate = validator.New(validator.WithRequiredStructEnabled())
func NewValidator() *XValidator {
	return &XValidator{Validator: validate}
}
func (v XValidator) Validate(data interface{}) []ErrorResponse {
    validationErrors := []ErrorResponse{}

    errs := validate.Struct(data)
    if errs != nil {
        for _, err := range errs.(validator.ValidationErrors) {
            // In this case data object is actually holding the User struct
            var elem ErrorResponse

            elem.FailedField = err.Field() // Export struct field name
            elem.Tag = err.Tag()           // Export struct tag
            elem.Value = err.Value()       // Export field value
            elem.Error = true

            validationErrors = append(validationErrors, elem)
        }
    }

    return validationErrors
}
