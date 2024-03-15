package dialog

import (
	"fmt"
)

const (
    TRANSFER_CONN_TYPE = "tcp"
    SERVER_CONN_TYPE = "tcp"
)

const (
    CLIENT = iota
    SERVER
)

const FIELD_PREFIX_LEN = 2
const OPTIMAL_FIELD_COUNT = 3

const (
    BLACKLIST = iota
    WHITELIST
)

const (
    ERROR = iota
    PING
    SIGN_UP
    UPDATE_IP
    GET_USER_DATA
    USER_DATA
    GET_LIST
    UPDATE_LIST
    TRANSFER_REQUEST
    TRANSFER_ACCEPTED
    TRANSFER
)

const PING_FIELD = ""

type MsgTypeErr struct {
    received, expected byte
    operation string
}

func (e *MsgTypeErr) Error() string {
    return fmt.Sprintf("error: Got message code %d instead of %d while in %s", e.received, e.expected, e.operation)
}

func (e *MsgTypeErr) Is(target error) bool {
    return target.Error() == e.Error()
}
