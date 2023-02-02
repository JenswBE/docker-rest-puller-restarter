package main

import (
	"fmt"

	"github.com/JenswBE/docker-rest-puller-restarter/openapi"
)

// APIError allows to bundle a status with the original error.
// This allows to fine-grained response codes at the API level.
type APIError struct {
	// HTTP status code
	Status int `json:"status"`

	// Original error
	Err error `json:"-"`

	// Error code
	Code string `json:"code"`

	// Human-readable description of the error
	Message string `json:"message"`

	// Optional - On which object to error occurred
	Instance string `json:"instance"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d - %s - %s - %s", e.Status, e.Message, e.Instance, e.Err.Error())
	}
	return fmt.Sprintf("%d - %s - %s", e.Status, e.Message, e.Instance)
}

// NewError returns a new APIError.
func NewError(status int, code openapi.APIErrorCode, instance string, err error) *APIError {
	return &APIError{
		Status:   status,
		Err:      err,
		Code:     string(code),
		Message:  translateCodeToMessage(code),
		Instance: instance,
	}
}

func translateCodeToMessage(code openapi.APIErrorCode) string {
	switch code {
	case openapi.APIERRORCODE_CONTAINER_NAME_NOT_ENABLED_FOR_API_KEY:
		return `Provided API key is not allowed to interact with the provided container name`
	case openapi.APIERRORCODE_INVALID_API_KEY:
		return `Provided API key is invalid`
	case openapi.APIERRORCODE_UNKNOWN_CONTAINER:
		return `There was no container found with the provided name`
	case openapi.APIERRORCODE_UNKNOWN_ERROR:
		return `An unknown error occurred`
	}
	return "" // Covered by exhaustive check
}
