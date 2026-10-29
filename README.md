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

## Development

Requires Go 1.27.

```sh
make test
make check
```
