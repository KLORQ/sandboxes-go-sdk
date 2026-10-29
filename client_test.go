package sandboxes

import (
	"context"
	"errors"
	"net"
	"testing"

	localv1 "github.com/KLORQ/sandboxes-api/gen/go/sandboxes/local/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestError(t *testing.T) {
	err := &Error{Code: "conflict", Message: "revision mismatch"}
	if got := err.Error(); got != "conflict: revision mismatch" {
		t.Fatalf("unexpected message: %s", got)
	}
}

func TestClientGetAndList(t *testing.T) {
	client, cleanup := newTestClient(t)
	defer cleanup()

	sb, err := client.Sandboxes().Get(context.Background(), "sb-1")
	if err != nil {
		t.Fatal(err)
	}
	if sb.ID != "sb-1" || sb.State != "running" {
		t.Fatalf("unexpected sandbox: %+v", sb)
	}

	page, err := client.Sandboxes().List(context.Background(), ListOptions{PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 || page.NextPageToken != "" {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestClientMapsNotFound(t *testing.T) {
	client, cleanup := newTestClient(t)
	defer cleanup()

	_, err := client.Sandboxes().Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	var typed *Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected typed error, got %T", err)
	}
	if typed.Code != "not_found" {
		t.Fatalf("code = %s", typed.Code)
	}
}

func TestOperationsWait(t *testing.T) {
	client, cleanup := newTestClient(t)
	defer cleanup()

	op, err := client.Operations().Wait(context.Background(), "op-1")
	if err != nil {
		t.Fatal(err)
	}
	if op.State != "done" {
		t.Fatalf("state = %s", op.State)
	}
}

func newTestClient(t *testing.T) (Client, func()) {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	localv1.RegisterSandboxServiceServer(server, fakeSandboxServer{})
	localv1.RegisterOperationServiceServer(server, fakeOperationServer{})
	go func() {
		_ = server.Serve(lis)
	}()
	client, err := New(Config{
		Address: "bufnet",
		Dialer: func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		},
	})
	if err != nil {
		server.Stop()
		t.Fatal(err)
	}
	return client, func() {
		_ = client.Close()
		server.Stop()
	}
}

type fakeSandboxServer struct {
	localv1.UnimplementedSandboxServiceServer
}

func (fakeSandboxServer) Get(_ context.Context, req *localv1.GetRequest) (*localv1.GetResponse, error) {
	if req.GetId() == "missing" {
		return nil, status.Error(codes.NotFound, "sandbox not found")
	}
	return &localv1.GetResponse{Sandbox: &localv1.Sandbox{
		Id:            "sb-1",
		Name:          "dev",
		State:         localv1.SandboxState_SANDBOX_STATE_RUNNING,
		Generation:    2,
		CreatedAtUnix: 100,
	}}, nil
}

func (fakeSandboxServer) List(context.Context, *localv1.ListRequest) (*localv1.ListResponse, error) {
	return &localv1.ListResponse{Items: []*localv1.Sandbox{{Id: "sb-1"}}}, nil
}

type fakeOperationServer struct {
	localv1.UnimplementedOperationServiceServer
}

func (fakeOperationServer) Get(context.Context, *localv1.OperationServiceGetRequest) (*localv1.OperationServiceGetResponse, error) {
	return &localv1.OperationServiceGetResponse{Operation: &localv1.Operation{
		Id:        "op-1",
		SandboxId: "sb-1",
		Kind:      localv1.OperationKind_OPERATION_KIND_CREATE,
		State:     localv1.OperationState_OPERATION_STATE_DONE,
	}}, nil
}
