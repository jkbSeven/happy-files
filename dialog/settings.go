package dialog

import (
	"errors"
	"fmt"
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

var defaultClient = map[string]any{
    "username": "",
    "email": "",
    "private-key": "",
    "server-ip": "localhost",
    "server-port": "13333",
    "server-public": "",
    "download-path": "~/Downloads"}

func msgBytes(code int, msg... string) []byte {
    return []byte(string(rune(code)) + "|" + strings.Join(msg, "|"))
}

func padWithZeros(field string, outLength int) ([]byte, error) {
    fieldLength := len(field)
    if outLength < fieldLength {
        return []byte{}, errors.New("String is longer than the desired length")
    }

    out := make([]byte, outLength)
    copy(out, field)

    return out, nil
}


