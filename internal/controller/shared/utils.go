package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// WriteErrorResponse writes an HTTP response which signify that an error has occurred.
func WriteErrorResponse(ginCtx *gin.Context, params WriteErrorResponseParams) {
	if params.Err != nil {
		params.Span.RecordError(params.Err)
		params.Span.SetStatus(codes.Error, params.Err.Error())
	}

	statusCode := params.StatusCode
	if statusCode == 0 {
		statusCode = http.StatusBadRequest
	}

	ginCtx.JSON(statusCode, params.Body)
}

type WriteErrorResponseParams struct {
	Body ErrorResponse
	Span trace.Span
	// Err is the error that will be recorded with the current span. If not provided, the error will not be recorded.
	Err error
	// StatusCode is an optional parameter to set the status code of the response.
	// If not provided, the default is http.StatusBadRequest.
	StatusCode int
}

// WriteBindingErrorResponse writes an HTTP response when a binding error (request body, query, ...) occurs.
func WriteBindingErrorResponse(ginCtx *gin.Context, span trace.Span, err error) {
	WriteErrorResponse(ginCtx, WriteErrorResponseParams{
		Body: ErrorResponse{
			Code:    CodeBindingRequestError,
			Message: err.Error(),
		},
		Span:       span,
		Err:        err,
		StatusCode: http.StatusBadRequest,
	})
}

// WriteUnknownErrorResponse writes an HTTP response when an unknown error occurs.
func WriteUnknownErrorResponse(ginCtx *gin.Context, span trace.Span, err error) {
	WriteErrorResponse(ginCtx, WriteErrorResponseParams{
		Body: ErrorResponse{
			Code:    CodeUnknown,
			Message: err.Error(),
		},
		Span:       span,
		Err:        err,
		StatusCode: http.StatusInternalServerError,
	})
}
