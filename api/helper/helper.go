package helper

import "github.com/go-playground/validator/v10"

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	JsonResponse := Response{
		Meta: meta,
		Data: data,
	}
	return JsonResponse
}
func FormatValidationError(err error) []string {
	var errors []string
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, e.Error())
		}
	} else {
		// Handle other types of errors here, if needed.
		// For example, if the error is not of type ValidationErrors,
		// you can add a default error message to the errors slice.
		errors = append(errors, "Validation error: "+err.Error())
	}
	return errors
}

