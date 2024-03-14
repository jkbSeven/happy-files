package dialog

import (
	"encoding/binary"
	"errors"
	"net"
	"syscall"
)

func bLength(length int) []byte {
    out := make([]byte, FIELD_PREFIX_LEN)
    binary.BigEndian.PutUint16(out, uint16(length))
    return out
}

func nLength(length []byte) int {
    return int(binary.BigEndian.Uint16(length))
}

func isAlive(conn net.Conn) bool {
    if conn == nil {
        return false
    } else if err := ping(conn); errors.Is(err, net.ErrClosed) || errors.Is(err, syscall.EPIPE) {
        return false
    }

    return true
}

func genMsg(code byte, fields... string) []byte {
    temp := []byte{}
    for _, v := range fields {
        length := bLength(len(v))
        temp = append(temp, length...)
        temp = append(temp, v...)
    }
    out := make([]byte, 1+ FIELD_PREFIX_LEN + len(temp))
    msgbLen := bLength(len(temp))
    out[0] = code
    out[1], out[2] = msgbLen[0], msgbLen[1]
    copy(out[3:], temp)

    return out
}

func readMsg(conn net.Conn) (byte, []byte, error) {
    msgData := make([]byte, 1 + FIELD_PREFIX_LEN)
    if _, err := conn.Read(msgData); err != nil {
        return 0, []byte{}, err
    }

    msgCode := msgData[0]
    msgLen := nLength(msgData[1:])

    if msgLen < 1 {
        return msgCode, []byte{}, nil
    }

    msg := make([]byte, nLength(msgData[1:]))

    if _, err := conn.Read(msg); err != nil {
        return 0, []byte{}, err
    }

    return msgCode, msg, nil
}

func groupMsg(msg []byte, msgLen int) [][]byte {
    msgFields := make([][]byte, 0, OPTIMAL_FIELD_COUNT)

    for read := 0; read < msgLen; {
        fieldLen := nLength(msg[:FIELD_PREFIX_LEN])
        msg = msg[FIELD_PREFIX_LEN:] // gets rid of FIELD_PREFIX
        if fieldLen < 1 {
            msgFields = append(msgFields, []byte{})
            read += FIELD_PREFIX_LEN
            continue
        }

        msgFields = append(msgFields, msg[:fieldLen])
        msg = msg[fieldLen:]
        read += FIELD_PREFIX_LEN + fieldLen
    }

    return msgFields
}

func rPadWithZeros(field string, outLength int) ([]byte, error) {
    fieldLen := len(field)
    if outLength < fieldLen {
        return []byte{}, errors.New("String is longer than the desired length")
    }

    out := make([]byte, outLength)
    copy(out, field)

    return out, nil
}

func ping(conn net.Conn) error {
    msg := genMsg(PING, PING_FIELD)
    if _, err := conn.Write(msg); err != nil {
        return err
    }

    msgCode, _, err := readMsg(conn)
    if err != nil {
        return err
    } else if msgCode != PING {
        return errors.New("Invalid response to ping")
    }

    return nil
}

func respondPing(conn net.Conn) error {
    msg := genMsg(PING, PING_FIELD)
    if _, err := conn.Write(msg); err != nil {
        return err
    }
    return nil
}
