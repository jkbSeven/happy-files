package dialog

import (
    "crypto/sha256"
    "crypto/rsa"
    "encoding/hex"
    // "fmt"
)

func hashPasswd(passwd string) string {
    sha := sha256.New()
    sha.Write([]byte(passwd))
    hpasswd := sha.Sum(nil)

    return hex.EncodeToString(hpasswd)
}

func GenerateKeys(key_size int) (*rsa.PrivateKey, error) {
    return nil, nil
}

func (client *Client) handshakeInit(username string) error {
    // client sends message (TRANSFER_REQUEST, username) and handles the rest of authentication
    return nil
}

func (client *Client) handshakeRespond(username string) error {
    // client got message (TRANSFER_REQUEST, username) and accepted it
    // this function handles authentication
    return nil
}
