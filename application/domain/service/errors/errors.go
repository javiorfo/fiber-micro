package errors

import (
	"net/http"

	"github.com/javiorfo/go-microservice-lib/response"
	"github.com/javiorfo/go-microservice-lib/response/backend"
	"go.opentelemetry.io/otel/trace"
)

func PermissionNotFound(span trace.Span) backend.Error {
	return response.NewResponseError(span,
		response.Error{
			HttpStatus: http.StatusBadRequest,
			Code:       response.ErrorCode("FIBER-MICRO-001"),
			Message:    response.Message("Permission not found!"),
		},
	)
}

func UserNotFound(span trace.Span) backend.Error {
	return response.NewResponseError(span,
		response.Error{
			HttpStatus: http.StatusBadRequest,
			Code:       response.ErrorCode("FIBER-MICRO-002"),
			Message:    response.Message("User not found!"),
		},
	)
}

func CredentialsError(span trace.Span) backend.Error {
	return response.NewResponseError(span,
		response.Error{
			HttpStatus: http.StatusUnauthorized,
			Code:       response.ErrorCode("FIBER-MICRO-003"),
			Message:    response.Message("Username or password incorrect!"),
		},
	)
}
