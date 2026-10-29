package sandboxes

// Public state values used by sandbox and operation records.
const (
	StatePending  = "pending"
	StateCreating = "creating"
	StateRunning  = "running"
	StateStopping = "stopping"
	StateStopped  = "stopped"
	StateDeleting = "deleting"
	StateDeleted  = "deleted"
	StateDone     = "done"
	StateFailed   = "failed"
	StateUnknown  = "unknown"
)
