package dmidecode

import (
	"fmt"
)

type ErrOpenStream struct {
	Err error
}

func (err ErrOpenStream) Error() string {
	return fmt.Sprintf("failed to open stream: %v", err.Err)
}

func (err ErrOpenStream) Unwrap() error {
	return err.Err
}

type ErrDecode struct {
	Err error
}

func (err ErrDecode) Error() string {
	return fmt.Sprintf("failed to decode structures: %v", err.Err)
}

func (err ErrDecode) Unwrap() error {
	return err.Err
}
