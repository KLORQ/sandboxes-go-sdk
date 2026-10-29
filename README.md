# sandboxes-go-sdk

Status: pre-alpha.

Official Go client for the local sandboxes daemon. The SDK is a client library,
not a second business-logic layer.

## Boundaries

- public domain types are not protobuf types;
- retries, idempotency, operation wait/watch, stream reconnect, pagination, and
  typed errors belong here;
- policy and security decisions belong to the daemon and gateway;
- the CLI and TUI consume this module instead of daemon internals.

## Usage

```go
client, err := sandboxes.New(sandboxes.Config{
	Address: "/path/to/sandboxes.sock",
})
if err != nil {
	return err
}
defer client.Close()

sb, err := client.Sandboxes().Get(ctx, "dev")
```

## Development

Requires Go 1.27.

```sh
make test
make check
make lint
```
