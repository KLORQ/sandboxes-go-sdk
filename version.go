package sandboxes

import "context"

// VersionInfo describes the daemon and API contract versions.
type VersionInfo struct {
	// Version is the daemon version.
	Version string
	// APIVersion is the local sandbox contract version implemented by the daemon.
	APIVersion string
}

// VersionService reports daemon and API contract versions.
type VersionService interface {
	Get(ctx context.Context) (VersionInfo, error)
}

// VersionServiceProvider is optionally implemented by root clients that expose
// daemon version inspection without changing the root Client interface.
type VersionServiceProvider interface {
	VersionService() VersionService
}

// Version returns the optional version service.
func Version(c Client) VersionService {
	if provider, ok := c.(VersionServiceProvider); ok {
		if service := provider.VersionService(); service != nil {
			return service
		}
	}
	return unsupportedVersionService{}
}

type unsupportedVersionService struct{}

func (unsupportedVersionService) Get(context.Context) (VersionInfo, error) {
	return VersionInfo{}, &Error{Code: errorCodeUnsupported, Message: "version inspection is not supported by this client"}
}
