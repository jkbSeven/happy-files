package dialog

import (
	"errors"
    "encoding/binary"
	"encoding/json"
	"os"
)

const TRANSFER_CONN_TYPE = "tcp"
const SERVER_CONN_TYPE = "tcp"

const (
    BLACKLIST = 0
    WHITELIST = 1
)

const (
    ERROR = 0
    SIGN_UP = 1
    UPDATE_IP = 2
    GET_USER_DATA = 3
    TRANSFER_REQUEST = 4
    GET_LIST = 5
    UPDATE_LIST = 6
)

var defaultClientConfig = map[string]any{
    "username": "",
    "email": "",
    "private-key": "",
    "server-ip": "localhost",
    "server-port": "13333",
    "server-public": "",
    "download-path": "~/Downloads",
}

var defaultServerConfig = map[string]any{
    "private-key": "",
    "db-path": "./hf.db",
}

func bLength(length int) []byte {
    out := make([]byte, 2)
    binary.BigEndian.PutUint16(out, uint16(length))
    return out
}

func genMsg(code byte, fields... string) []byte {
    temp := []byte{}
    for _, v := range fields {
        length := bLength(len(v))
        temp = append(temp, length...)
        temp = append(temp, v...)
    }
    out := make([]byte, 3 + len(temp))
    msgbLen := bLength(len(temp))
    out[0], out[1] = msgbLen[0], msgbLen[1]
    out[2] = code
    copy(out[3:], temp)

    return out
}

func rPadWithZeros(field string, outLength int) ([]byte, error) {
    fieldLength := len(field)
    if outLength < fieldLength {
        return []byte{}, errors.New("String is longer than the desired length")
    }

    out := make([]byte, outLength)
    copy(out, field)

    return out, nil
}

func readConfig(configPath string) (map[string]any, error) {
    var config map[string]any

    data, err := os.ReadFile(configPath)
    if err != nil {
        return config, err
    }

    err = json.Unmarshal(data, &config)
    if err != nil {
        return config, err
    }

    return config, nil
}

func writeConfig(configPath string, changes map[string]any) error {
    config, err := readConfig(configPath)
    if err != nil {
        return err
    }

    for k, v := range changes {
        config[k] = v
    }

    
    data, err := json.Marshal(config)
    if err != nil {
        return err
    }

    err = os.WriteFile(configPath, data, 0644)
    if err != nil {
        return err
    }

    return nil
}

