package error

import "strings"

//HTTP Error messages
const (
	InternalServerError        string = "Internal Server Error"
	BadRequestError                   = "Bad Request Error"
	UnauthorizedError                 = "Unauthorized Error"
	NotFoundError                     = "Not Found Error"
	ThresholdError                    = "Threshold Error"
	ForbiddenError                    = "Forbidden Error"
	AppError                          = "App Error"
	HealthError                       = "Health Error"
)

//Error codes
const (
	//Auth
	InvalidTokenErrorCode = "invalid-token"
	UnauthorizedErrorCode = "unauthorized"
	//general
	ForbiddenErrorCode              = "forbidden-operation"
	InvalidOperationErrorCode       = "invalid-operation"
	RequestValidationErrorCode      = "request-validation-error"
	ValidationErrorCode             = "validation-error"
	InvalidExpressionErrorCode      = "invalid-expression"
	InternalServerErrorCode         = "internal-server-error"
	InvalidPayloadErrorCode         = "invalid-payload"
	EvaluationErrorCode             = "evaluation-error"
	SchemaValidationFailedErrorCode = "schema-validation-failed"
	MarshallingErrorCode            = "marshalling-failed"
	UnMarshallingErrorCode          = "unmarshalling-failed"
	GlobalPanicErrorCode            = "global-panic"
)

//Error messages
const (
	InvalidRequestPayload = "Invalid Request Payload"
	InvalidAuthHeader     = "Invalid Authorization Header"
	InvalidToken          = "Invalid Auth token"
	InvalidFilter         = "Invalid Filter"
	InvalidProjection     = "Invalid Projection"
)

//used to append to error code, example - code := {component}/{error-code}
const (
	//Auth component
	Auth string = "auth"
	//App component
	App string = "app"
	//Health component
	Health string = "health"
)

//GenerateErrorCode ..
func GenerateErrorCode(args ...string) string {
	return strings.ToLower(strings.Join(args, "/"))
}
