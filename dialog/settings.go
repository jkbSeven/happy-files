package dialog

import (
    "strings"
)

const TRANSFER_CONN_TYPE = "tcp"
const SERVER_CONN_TYPE = "tcp"
const (
    ERROR = 0
    SIGN_UP = 1
    UPDATE = 2
    SYNC = 3
    GET_USER_DATA = 4
    TRANSFER_REQUEST = 5
)

func msgBytes(code int, msg... string) []byte {
    return []byte(string(rune(code)) + "|" + strings.Join(msg, "|"))
}

