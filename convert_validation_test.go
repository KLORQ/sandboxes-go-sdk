package sandboxes

import (
	"testing"

	localv1 "github.com/KLORQ/sandboxes-api/gen/go/sandboxes/local/v1"
)

func TestSandboxConversionRejectsInvalidNetworkFields(t *testing.T) {
	for name, input := range map[string]*localv1.Sandbox{
		"host port":  {Ports: []*localv1.Port{{Host: 65536}}},
		"guest port": {Ports: []*localv1.Port{{Guest: 65536}}},
		"allow rule": {Allow: []*localv1.NetworkRule{{Port: 65536}}},
		"deny rule":  {Deny: []*localv1.NetworkRule{{Port: 65536}}},
		"nil rule":   {Deny: []*localv1.NetworkRule{nil}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := fromProtoSandbox(input); err == nil {
				t.Fatal("invalid daemon response accepted")
			}
		})
	}
	valid := &localv1.Sandbox{Ports: []*localv1.Port{{Host: 65535, Guest: 443}}, Deny: []*localv1.NetworkRule{{Port: 65535}}}
	if _, err := fromProtoSandbox(valid); err != nil {
		t.Fatalf("valid daemon response rejected: %v", err)
	}
}
