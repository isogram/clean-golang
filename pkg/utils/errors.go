package utils

import (
	"errors"
)

const (
	// ErrorDataNotFound error message when data doesn't exist
	ErrorDataNotFound = "data %s not found"
	// ErrorParameterInvalid error message for parameter is invalid
	ErrorParameterInvalid = "%s parameter is invalid"
	// ErrorParameterRequired error message for parameter is missing
	ErrorParameterRequired = "%s parameter is required"
	// ErrorParameterLength error message for parameter length is invalid
	ErrorParameterLength = "length of %s parameter exceeds the limit %d"
	// ErrorUnauthorized error message for unauthorized user
	ErrorUnauthorized = "you are not authorized"
	// ErrorPayloadInvalid error message when payload is invalid
	ErrorPayloadInvalid = "payload is invalid"
	// ErrorProcessingRequest error message when processing request
	ErrorProcessingRequest = "failed to process request"
	// ErrorMissingEnvVariable error message when env not set
	ErrorMissingEnvVariable = "you need to specify %v in the environment variable"
)

// ErrMapStringMessage map string error message for raven
type ErrMapStringMessage map[string]string

var (
	// ErrorFormatDate error
	ErrorFormatDate = errors.New("format date should be yyyy-mm-dd")
)
