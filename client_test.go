package sandboxes

import "testing"

func TestError(t *testing.T) {
	err := &Error{Code: "conflict", Message: "revision mismatch"}
	if got := err.Error(); got != "conflict: revision mismatch" {
		t.Fatalf("unexpected message: %s", got)
	}
}
