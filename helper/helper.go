package helper

import "github.com/go-playground/validator/v10"

type Response struct {
	MetaMessage Meta        `json:"meta"`
	Info        interface{} `json:"info_data"`
	Data        interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func ApiResponse(message string, code int, status string, infoData, data interface{}) Response {
	metaData := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		MetaMessage: metaData,
		Info:        infoData,
		Data:        data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}
