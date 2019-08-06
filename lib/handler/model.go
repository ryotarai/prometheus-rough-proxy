package handler

import "github.com/prometheus/common/model"

type status string

const (
	statusSuccess status = "success"
	statusError   status = "error"
)

type errorType string

const (
	errorNone        errorType = ""
	errorTimeout     errorType = "timeout"
	errorCanceled    errorType = "canceled"
	errorExec        errorType = "execution"
	errorBadData     errorType = "bad_data"
	errorInternal    errorType = "internal"
	errorUnavailable errorType = "unavailable"
	errorNotFound    errorType = "not_found"
)

type apiResponse struct {
	Status    status      `json:"status"`
	Data      interface{} `json:"data"`
	ErrorType errorType   `json:"errorType,omitEmpty"`
	Error     string      `json:"error,omitEmpty"`
}

type queryData struct {
	ResultType model.ValueType `json:"resultType"`
	Result     model.Matrix    `json:"result"`
}
