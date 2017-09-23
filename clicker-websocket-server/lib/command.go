package cc_websocket_server

import (
	"fmt"
)

type ClientConnectionStatus int

const (
	NEW_CONNECTION_STATUS ClientConnectionStatus = iota
	RUNNING_CONNECTION_STATUS
	CONNECTION_STATUS_EOF
)

type CommandResponseErrorType int

const (
	ERROR_TYPE_SUCCESS CommandResponseErrorType = iota
	ERROR_TYPE_UNAUTHORIZED
	ERROR_TYPE_INVALID_REQUEST
	ERROR_TYPE_SERVICE_ERROR
)

type MethodType string

const (
	METHOD_TYPE_POST   MethodType = "POST"
	METHOD_TYPE_LIST   MethodType = "LIST"
	METHOD_TYPE_DELETE MethodType = "DELETE"
)

var METHOD_TYPE_LOOKUP map[MethodType]bool = map[MethodType]bool{
	METHOD_TYPE_POST:   true,
	METHOD_TYPE_LIST:   true,
	METHOD_TYPE_DELETE: true,
}

type ObjectType string

const (
	OBJECT_TYPE_GAME     ObjectType = "game"
	OBJECT_TYPE_UPGRADE  ObjectType = "upgrade"
	OBJECT_TYPE_BUILDING ObjectType = "building"
	OBJECT_TYPE_COOKIE   ObjectType = "cookie"
)

var OBJECT_TYPE_LOOKUP map[ObjectType]bool = map[ObjectType]bool{
	OBJECT_TYPE_GAME:     true,
	OBJECT_TYPE_UPGRADE:  true,
	OBJECT_TYPE_BUILDING: true,
	OBJECT_TYPE_COOKIE:   true,
}

type CommandRequest struct {
	Object ObjectType `json:"type"`
	ID     string     `json:"id"`
	Method MethodType `json:"method"`
}

type CommandFunction func(*Client, *CommandRequest, *CommandResponse)
type SupportedCommand struct {
	object ObjectType
	hasID  bool
	method MethodType
}

func NullFunction(c *Client, req *CommandRequest, resp *CommandResponse) {}

var COMMAND_DISPATCH_TABLE map[SupportedCommand]CommandFunction = map[SupportedCommand]CommandFunction{
	SupportedCommand{
		object: OBJECT_TYPE_GAME,
		hasID:  false,
		method: METHOD_TYPE_POST,
	}: NullFunction,
	SupportedCommand{
		object: OBJECT_TYPE_GAME,
		hasID:  true,
		method: METHOD_TYPE_DELETE,
	}: NullFunction,
	SupportedCommand{
		object: OBJECT_TYPE_UPGRADE,
		hasID:  false,
		method: METHOD_TYPE_LIST,
	}: NullFunction,
	SupportedCommand{
		object: OBJECT_TYPE_UPGRADE,
		hasID:  true,
		method: METHOD_TYPE_POST,
	}: NullFunction,
	SupportedCommand{
		object: OBJECT_TYPE_BUILDING,
		hasID:  false,
		method: METHOD_TYPE_LIST,
	}: NullFunction,
	SupportedCommand{
		object: OBJECT_TYPE_BUILDING,
		hasID:  true,
		method: METHOD_TYPE_POST,
	}: NullFunction,
	SupportedCommand{
		object: OBJECT_TYPE_COOKIE,
		hasID:  false,
		method: METHOD_TYPE_POST,
	}: NullFunction,
}

func (self *CommandRequest) validate(command_error *CommandError) {
	_, valid_object_type := OBJECT_TYPE_LOOKUP[self.Object]
	if !valid_object_type {
		command_error.ErrorCode = ERROR_TYPE_INVALID_REQUEST
		command_error.ErrorMessage = fmt.Sprintf("Invalid object type '%s'", self.Object)
	}

	_, valid_method_type := METHOD_TYPE_LOOKUP[self.Method]
	if !valid_method_type {
		command_error.ErrorCode = ERROR_TYPE_INVALID_REQUEST
		command_error.ErrorMessage = fmt.Sprintf("Unsupported command type '%s'", self.Method)
	}

	_, supported_method_call := COMMAND_DISPATCH_TABLE[SupportedCommand{
		object: self.Object,
		hasID:  self.ID != "",
		method: self.Method,
	}]
	if !supported_method_call {
		command_error.ErrorCode = ERROR_TYPE_INVALID_REQUEST
		command_error.ErrorMessage = fmt.Sprintf("Unsupported command type '%s' for object '%s' and ID '%s'", self.Method, self.Object, self.ID)
	}
}

type CommandError struct {
	ErrorCode    CommandResponseErrorType `json:"error_code"`
	ErrorMessage string                   `json:"error_message"`
}

type CommandResponse struct {
	Content CommandRequest `json:"content"`
	Error   CommandError   `json:"error"`
}
