/*
Docker Pull-Restarter

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// APIErrorCode - CONTAINER_NAME_NOT_ENABLED_FOR_API_KEY: Provided API key is not allowed to interact with the provided container name - INVALID_API_KEY: Provided API key is invalid - UNKNOWN_CONTAINER: There was no container found with the provided name - UNKNOWN_ERROR: An unknown error occurred 
type APIErrorCode string

// List of APIErrorCode
const (
	APIERRORCODE_CONTAINER_NAME_NOT_ENABLED_FOR_API_KEY APIErrorCode = "CONTAINER_NAME_NOT_ENABLED_FOR_API_KEY"
	APIERRORCODE_INVALID_API_KEY APIErrorCode = "INVALID_API_KEY"
	APIERRORCODE_UNKNOWN_CONTAINER APIErrorCode = "UNKNOWN_CONTAINER"
	APIERRORCODE_UNKNOWN_ERROR APIErrorCode = "UNKNOWN_ERROR"
)

// All allowed values of APIErrorCode enum
var AllowedAPIErrorCodeEnumValues = []APIErrorCode{
	"CONTAINER_NAME_NOT_ENABLED_FOR_API_KEY",
	"INVALID_API_KEY",
	"UNKNOWN_CONTAINER",
	"UNKNOWN_ERROR",
}

func (v *APIErrorCode) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := APIErrorCode(value)
	for _, existing := range AllowedAPIErrorCodeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid APIErrorCode", value)
}

// NewAPIErrorCodeFromValue returns a pointer to a valid APIErrorCode
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewAPIErrorCodeFromValue(v string) (*APIErrorCode, error) {
	ev := APIErrorCode(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for APIErrorCode: valid values are %v", v, AllowedAPIErrorCodeEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v APIErrorCode) IsValid() bool {
	for _, existing := range AllowedAPIErrorCodeEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to APIErrorCode value
func (v APIErrorCode) Ptr() *APIErrorCode {
	return &v
}

type NullableAPIErrorCode struct {
	value *APIErrorCode
	isSet bool
}

func (v NullableAPIErrorCode) Get() *APIErrorCode {
	return v.value
}

func (v *NullableAPIErrorCode) Set(val *APIErrorCode) {
	v.value = val
	v.isSet = true
}

func (v NullableAPIErrorCode) IsSet() bool {
	return v.isSet
}

func (v *NullableAPIErrorCode) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAPIErrorCode(val *APIErrorCode) *NullableAPIErrorCode {
	return &NullableAPIErrorCode{value: val, isSet: true}
}

func (v NullableAPIErrorCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAPIErrorCode) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

