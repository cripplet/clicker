package cc_websocket_server

import (
	"testing"
)

func TestInvalidRequestInvalidObject(t *testing.T) {
	e := CommandError{}
	c := CommandRequest{
		Object: "invalid-object",
	}
	(&c).validate(&e)
	if e.ErrorCode != ERROR_TYPE_INVALID_REQUEST {
		t.Errorf("Unexpected error: %d != %d", e.ErrorCode, ERROR_TYPE_INVALID_REQUEST)
	}
}

func TestInvalidRequestInvalidMethod(t *testing.T) {
	e := CommandError{}
	c := CommandRequest{
		Object: "game",
		Method: "INVALID_METHOD",
	}
	(&c).validate(&e)
	if e.ErrorCode != ERROR_TYPE_INVALID_REQUEST {
		t.Errorf("Unexpected error: %d != %d", e.ErrorCode, ERROR_TYPE_INVALID_REQUEST)
	}
}

func TestUnsuppportedCommand(t *testing.T) {
	e := CommandError{}
	c := CommandRequest{
		Object: "game",
		Method: "LIST",
	}
	(&c).validate(&e)
	if e.ErrorCode != ERROR_TYPE_INVALID_REQUEST {
		t.Errorf("Unexpected error: %d != %d", e.ErrorCode, ERROR_TYPE_INVALID_REQUEST)
	}
}

func TestSupportedCommand(t *testing.T) {
	e := CommandError{}
	c := CommandRequest{
		Object: "game",
		Method: "POST",
	}
	(&c).validate(&e)
	if e.ErrorCode != ERROR_TYPE_SUCCESS {
		t.Errorf("Unexpected error: %d (%s)", e.ErrorCode, e.ErrorMessage)
	}
}
