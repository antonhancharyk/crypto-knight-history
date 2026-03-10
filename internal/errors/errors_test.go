package errors

import (
	"errors"
	"testing"
)

func TestSentinels(t *testing.T) {
	if !errors.Is(ErrNotFound, ErrNotFound) {
		t.Error("errors.Is(ErrNotFound, ErrNotFound) should be true")
	}
	if !errors.Is(ErrBadRequest, ErrBadRequest) {
		t.Error("errors.Is(ErrBadRequest, ErrBadRequest) should be true")
	}
	wrapped := errors.Join(ErrNotFound, errors.New("detail"))
	if !errors.Is(wrapped, ErrNotFound) {
		t.Error("wrapped error should be ErrNotFound")
	}
}
