package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ResponseFormat: struct hold api response format
type ResponseFormat struct {
	Status  int               `json:"status,omitempty"`
	Message interface{}       `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Total   interface{}       `json:"total,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// ErrorValidator: struct hold element of validator error
type ErrorValidator struct {
	FailedField string
	Tag         string
	Value       string
}

// SetData: method to set data to response if any
func (r *ResponseFormat) SetData(data interface{}, total interface{}) *ResponseFormat {
	r.Status = http.StatusOK
	r.Message = "success"
	r.Data = data
	r.Total = total
	r.Errors = nil

	return r
}

// SetError: method to set error to response if any
func (r *ResponseFormat) SetError(err error) *ResponseFormat {
	r.Status = http.StatusBadRequest
	r.Message = "Fail"
	r.Errors = map[string]string{"Errors": err.Error()}
	r.Data = nil

	if httpErr, ok := err.(*fiber.Error); ok {
		r.Status = httpErr.Code
	} else if validationErr, ok := err.(validator.ValidationErrors); ok {
		var ErrorResponse = make(map[string]string)
		for _, err := range validationErr {
			var element ErrorValidator
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			ErrorResponse[element.FailedField] = element.Tag
		}
		r.Status = http.StatusUnprocessableEntity
		r.Errors = ErrorResponse
	}

	return r
}
