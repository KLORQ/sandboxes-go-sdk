package sandboxes

import (
	"context"
	"fmt"
	"time"

	localv1 "github.com/KLORQ/sandboxes-api/gen/go/sandboxes/local/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	conn        *grpc.ClientConn
	sandboxes   localv1.SandboxServiceClient
	operations  localv1.OperationServiceClient
	callTimeout time.Duration
}

// New creates a lazy local daemon client.
func New(cfg Config) (Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("sandbox daemon address is required")
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if cfg.Dialer != nil {
		opts = append(opts, grpc.WithContextDialer(cfg.Dialer))
	}
	conn, err := grpc.NewClient("unix:"+cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}
	return &client{
		conn:        conn,
		sandboxes:   localv1.NewSandboxServiceClient(conn),
		operations:  localv1.NewOperationServiceClient(conn),
		callTimeout: cfg.CallTimeout,
	}, nil
}

func (c *client) Sandboxes() SandboxService { return sandboxClient{c: c} }

func (c *client) Operations() OperationService { return operationClient{c: c} }

func (c *client) Close() error { return c.conn.Close() }

func (c *client) context(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.callTimeout <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, c.callTimeout)
}

type sandboxClient struct{ c *client }

func (s sandboxClient) Create(ctx context.Context, req CreateRequest) (*Operation, error) {
	ctx, cancel := s.c.context(ctx)
	defer cancel()
	resp, err := s.c.sandboxes.Create(ctx, &localv1.CreateRequest{
		Name:           req.Name,
		Image:          req.Image,
		Workspace:      req.Workspace,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return fromProtoOperation(resp.GetOperation()), nil
}

func (s sandboxClient) Get(ctx context.Context, id string) (*Sandbox, error) {
	ctx, cancel := s.c.context(ctx)
	defer cancel()
	resp, err := s.c.sandboxes.Get(ctx, &localv1.GetRequest{Id: id})
	if err != nil {
		return nil, mapError(err)
	}
	return fromProtoSandbox(resp.GetSandbox()), nil
}

func (s sandboxClient) List(ctx context.Context, opts ListOptions) (*Page[Sandbox], error) {
	if opts.PageSize < 0 || opts.PageSize > 1<<31-1 {
		return nil, &Error{Code: "invalid_argument", Message: "page size out of range"}
	}
	ctx, cancel := s.c.context(ctx)
	defer cancel()
	resp, err := s.c.sandboxes.List(ctx, &localv1.ListRequest{
		PageSize:  int32(opts.PageSize),
		PageToken: opts.PageToken,
	})
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]Sandbox, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		if sb := fromProtoSandbox(item); sb != nil {
			items = append(items, *sb)
		}
	}
	return &Page[Sandbox]{Items: items, NextPageToken: resp.GetNextPageToken()}, nil
}

func (s sandboxClient) Start(ctx context.Context, id string) (*Operation, error) {
	ctx, cancel := s.c.context(ctx)
	defer cancel()
	resp, err := s.c.sandboxes.Start(ctx, &localv1.StartRequest{Id: id})
	if err != nil {
		return nil, mapError(err)
	}
	return fromProtoOperation(resp.GetOperation()), nil
}

func (s sandboxClient) Stop(ctx context.Context, id string, opts StopOptions) (*Operation, error) {
	ctx, cancel := s.c.context(ctx)
	defer cancel()
	resp, err := s.c.sandboxes.Stop(ctx, &localv1.StopRequest{
		Id:           id,
		GraceSeconds: int64(opts.Grace / time.Second),
	})
	if err != nil {
		return nil, mapError(err)
	}
	return fromProtoOperation(resp.GetOperation()), nil
}

func (s sandboxClient) Delete(ctx context.Context, id string, opts DeleteOptions) (*Operation, error) {
	ctx, cancel := s.c.context(ctx)
	defer cancel()
	resp, err := s.c.sandboxes.Delete(ctx, &localv1.DeleteRequest{Id: id, Force: opts.Force})
	if err != nil {
		return nil, mapError(err)
	}
	return fromProtoOperation(resp.GetOperation()), nil
}

func (s sandboxClient) Watch(ctx context.Context, opts WatchOptions) (EventStream, error) {
	stream, err := s.c.sandboxes.Watch(ctx, &localv1.WatchRequest{
		SandboxId:     opts.SandboxID,
		AfterRevision: opts.AfterRevision,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &eventStream{stream: stream}, nil
}

type eventStream struct {
	stream grpc.ServerStreamingClient[localv1.WatchResponse]
}

func (e *eventStream) Recv() (*Event, error) {
	resp, err := e.stream.Recv()
	if err != nil {
		return nil, mapError(err)
	}
	return fromProtoEvent(resp), nil
}

func (e *eventStream) Close() error {
	return e.stream.CloseSend()
}

type operationClient struct{ c *client }

func (o operationClient) Get(ctx context.Context, id string) (*Operation, error) {
	ctx, cancel := o.c.context(ctx)
	defer cancel()
	resp, err := o.c.operations.Get(ctx, &localv1.OperationServiceGetRequest{Id: id})
	if err != nil {
		return nil, mapError(err)
	}
	return fromProtoOperation(resp.GetOperation()), nil
}

func (o operationClient) Wait(ctx context.Context, id string) (*Operation, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		op, err := o.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if op.State == "done" || op.State == "failed" {
			return op, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}
