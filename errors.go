package sandboxes

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	var typed *Error
	if errors.As(err, &typed) {
		return typed
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}

	code := "internal"
	retryable := false
	switch status.Code(err) {
	case codes.NotFound:
		code = "not_found"
	case codes.AlreadyExists:
		code = "already_exists"
	case codes.Aborted, codes.FailedPrecondition:
		code = "conflict"
	case codes.InvalidArgument:
		code = "invalid_argument"
	case codes.PermissionDenied, codes.Unauthenticated:
		code = "permission_denied"
	case codes.Unavailable:
		code, retryable = "unavailable", true
	case codes.DeadlineExceeded:
		code, retryable = "deadline_exceeded", true
	case codes.ResourceExhausted:
		code, retryable = "resource_exhausted", true
	}
	return &Error{Code: code, Message: status.Convert(err).Message(), Retryable: retryable}
}
