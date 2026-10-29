package sandboxes

import (
	"time"

	localv1 "github.com/KLORQ/sandboxes-api/gen/go/sandboxes/local/v1"
)

func fromProtoSandbox(in *localv1.Sandbox) *Sandbox {
	if in == nil {
		return nil
	}
	return &Sandbox{
		ID:         in.GetId(),
		Name:       in.GetName(),
		State:      sandboxState(in.GetState()),
		Generation: in.GetGeneration(),
		Image:      in.GetImage(),
		Workspace:  in.GetWorkspace(),
		Error:      in.GetError(),
		CreatedAt:  unixTime(in.GetCreatedAtUnix()),
		UpdatedAt:  unixTime(in.GetUpdatedAtUnix()),
	}
}

func fromProtoOperation(in *localv1.Operation) *Operation {
	if in == nil {
		return nil
	}
	return &Operation{
		ID:         in.GetId(),
		SandboxID:  in.GetSandboxId(),
		Kind:       operationKind(in.GetKind()),
		State:      operationState(in.GetState()),
		Generation: in.GetGeneration(),
		Error:      in.GetError(),
		CreatedAt:  unixTime(in.GetCreatedAtUnix()),
		UpdatedAt:  unixTime(in.GetUpdatedAtUnix()),
	}
}

func fromProtoEvent(in *localv1.WatchResponse) *Event {
	if in == nil {
		return nil
	}
	return &Event{
		Type:      in.GetType(),
		Sandbox:   fromProtoSandbox(in.GetSandbox()),
		Operation: fromProtoOperation(in.GetOperation()),
		Revision:  in.GetRevision(),
	}
}

func unixTime(seconds int64) time.Time {
	if seconds == 0 {
		return time.Time{}
	}
	return time.Unix(seconds, 0).UTC()
}

func sandboxState(v localv1.SandboxState) string {
	switch v {
	case localv1.SandboxState_SANDBOX_STATE_PENDING:
		return StatePending
	case localv1.SandboxState_SANDBOX_STATE_CREATING:
		return StateCreating
	case localv1.SandboxState_SANDBOX_STATE_RUNNING:
		return StateRunning
	case localv1.SandboxState_SANDBOX_STATE_STOPPING:
		return StateStopping
	case localv1.SandboxState_SANDBOX_STATE_STOPPED:
		return StateStopped
	case localv1.SandboxState_SANDBOX_STATE_DELETING:
		return StateDeleting
	case localv1.SandboxState_SANDBOX_STATE_DELETED:
		return StateDeleted
	case localv1.SandboxState_SANDBOX_STATE_FAILED:
		return StateFailed
	default:
		return StateUnknown
	}
}

func operationKind(v localv1.OperationKind) string {
	switch v {
	case localv1.OperationKind_OPERATION_KIND_CREATE:
		return "create"
	case localv1.OperationKind_OPERATION_KIND_START:
		return "start"
	case localv1.OperationKind_OPERATION_KIND_STOP:
		return "stop"
	case localv1.OperationKind_OPERATION_KIND_DELETE:
		return "delete"
	default:
		return StateUnknown
	}
}

func operationState(v localv1.OperationState) string {
	switch v {
	case localv1.OperationState_OPERATION_STATE_PENDING:
		return StatePending
	case localv1.OperationState_OPERATION_STATE_RUNNING:
		return StateRunning
	case localv1.OperationState_OPERATION_STATE_DONE:
		return StateDone
	case localv1.OperationState_OPERATION_STATE_FAILED:
		return StateFailed
	default:
		return StateUnknown
	}
}
