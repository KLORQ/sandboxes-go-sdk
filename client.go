// Package sandboxes provides the public Go client for the local daemon.
package sandboxes

import "context"

// Client is the root SDK client.
type Client interface {
	Sandboxes() SandboxService
	Operations() OperationService
	Close() error
}

// SandboxService manages sandbox lifecycle.
type SandboxService interface {
	Get(ctx context.Context, id string) (*Sandbox, error)
	List(ctx context.Context, opts ListOptions) (*Page[Sandbox], error)
}

// OperationService waits for and observes long-running operations.
type OperationService interface {
	Get(ctx context.Context, id string) (*Operation, error)
	Wait(ctx context.Context, id string) (*Operation, error)
}

// Sandbox is the public sandbox representation.
type Sandbox struct {
	ID         string
	Name       string
	State      string
	Generation int64
}

// Operation is the public operation representation.
type Operation struct {
	ID        string
	Kind      string
	State     string
	Retryable bool
}

// ListOptions controls pagination.
type ListOptions struct {
	PageSize  int
	PageToken string
}

// Page is one page of results.
type Page[T any] struct {
	Items         []T
	NextPageToken string
}

// Error is a typed SDK error.
type Error struct {
	Code      string
	Message   string
	Retryable bool
	Operation string
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }
