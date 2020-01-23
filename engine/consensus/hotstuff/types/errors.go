package types

import "fmt"

type ErrorFinalizationFatal struct {
}

func (e *ErrorFinalizationFatal) Error() string {
	panic("implement me")
}

type ErrorConfiguration struct {
	Msg string
}

func (e *ErrorConfiguration) Error() string {
	return e.Msg
}

type ErrorInvalidTimeout struct {
	Timeout *Timeout
	CurrentView uint64
	CurrentMode TimeoutMode
}

func (e *ErrorInvalidTimeout) Error() string {
	return fmt.Sprintf(
		"received timeout (view, mode) (%d, %s) but current state is (%d, %s)",
		e.Timeout.View, e.Timeout.Mode.String(), e.CurrentView, e.CurrentMode.String(),
		)
}

type ErrorConflictingQCs struct {
	View uint64
	Qcs []*QuorumCertificate
}

func (e *ErrorConflictingQCs) Error() string {
	return fmt.Sprintf("%d conflicting QCs at view %d", e.View, len(e.Qcs))
}

type ErrorMissingBlock struct {
	View uint64
	BlockID []byte
}

func (e *ErrorMissingBlock) Error() string {
	return fmt.Sprintf("missing Block at view %d with ID %s", e.View, string(e.BlockID))
}

type ErrorInvalidBlock struct {
	View uint64
	BlockID []byte
	Msg string
}

func (e *ErrorInvalidBlock) Error() string {
	return fmt.Sprintf("invalid block (view %d; ID %s): %s", e.View, string(e.BlockID), e.Msg)
}
