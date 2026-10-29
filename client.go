// Package sandboxes provides the public Go client for the local daemon.
package sandboxes

import (
	"context"
	"net"
	"time"
)

// Config controls the local daemon connection.
type Config struct {
	// Address is the Unix socket path of the local daemon.
	Address string
	// Dialer overrides the default Unix socket dialer.
	Dialer func(ctx context.Context, address string) (net.Conn, error)
	// CallTimeout bounds unary calls. Zero uses the context deadline.
	CallTimeout time.Duration
}

// Client is the root SDK client.
type Client interface {
	Sandboxes() SandboxService
	Operations() OperationService
	Close() error
}

// SandboxService manages sandbox lifecycle.
type SandboxService interface {
	Create(ctx context.Context, req CreateRequest) (*Operation, error)
	Get(ctx context.Context, id string) (*Sandbox, error)
	List(ctx context.Context, opts ListOptions) (*Page[Sandbox], error)
	Start(ctx context.Context, id string) (*Operation, error)
	Stop(ctx context.Context, id string, opts StopOptions) (*Operation, error)
	Delete(ctx context.Context, id string, opts DeleteOptions) (*Operation, error)
}

// OperationService observes long-running operations.
type OperationService interface {
	Get(ctx context.Context, id string) (*Operation, error)
	Wait(ctx context.Context, id string) (*Operation, error)
}

// CreateRequest describes a sandbox create request.
type CreateRequest struct {
	Name           string
	Image          string
	Workspace      string
	IdempotencyKey string
}

// StopOptions controls graceful shutdown.
type StopOptions struct {
	Grace time.Duration
}

// DeleteOptions controls deletion.
type DeleteOptions struct {
	Force bool
}

// Sandbox is the public sandbox representation.
type Sandbox struct {
	ID         string
	Name       string
	State      string
	Generation int64
	Image      string
	Workspace  string
	Error      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Operation is the public operation representation.
type Operation struct {
	ID         string
	SandboxID  string
	Kind       string
	State      string
	Generation int64
	Error      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
