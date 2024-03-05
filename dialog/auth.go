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
