package sandboxes

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errorCodeInvalidArgument = "invalid_argument"
const errorCodeUnsupported = "unsupported"
const errorCodeInternal = "internal"

func mapError(err error) error {
	if err == nil {
		return nil
	}
	var typed *Error
	if errors.As(err, &typed) {
		return typed
	}
	if errors.Is(err, io.EOF) {
		return io.EOF
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}

	code := errorCodeInternal
	retryable := false
	switch status.Code(err) {
	case codes.NotFound:
		code = "not_found"
	case codes.AlreadyExists:
		code = "already_exists"
	case codes.Aborted, codes.FailedPrecondition:
		code = "conflict"
	case codes.InvalidArgument:
		code = errorCodeInvalidArgument
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
