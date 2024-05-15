package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
)

type ClientConfig struct {
    Username string
    Email string
    PrivateKey *rsa.PrivateKey
    PublicKey *rsa.PublicKey
    StoreHistory bool
    DownloadPath string
    ServerIP string
    ServerPublicKey *rsa.PublicKey
}

type ServerConfig struct {
    Port string
    PrivateKey *rsa.PrivateKey
    PublicKey *rsa.PublicKey
    DatabasePath string
}

func ReadClientConfig(path string) (ClientConfig, error) {
    var clientConfig ClientConfig
    
    config, err := readConfig(path)
    if err != nil {
        return clientConfig, err
    }
    
    clientConfig.Username = config["Username"].(string)
    clientConfig.Email = config["Email"].(string)

    privateKey, err := parsePrivateKey(config["PrivateKey"].(string))
    if err != nil {
        panic(err)
    }
    clientConfig.PrivateKey = privateKey

    publicKey, err := parsePublicKey(config["PublicKey"].(string))
    if err != nil {
        panic(err)
    }
    clientConfig.PublicKey = publicKey

    clientConfig.StoreHistory = config["StoreHistory"].(bool)
    clientConfig.DownloadPath = config["DownloadPath"].(string)
    clientConfig.ServerIP = config["ServerIP"].(string)

    publicKey, err = parsePublicKey(config["ServerPublicKey"].(string))
    if err != nil {
        panic(err)
    }
    clientConfig.ServerPublicKey = publicKey

    return clientConfig, nil
}

func ReadServerConfig(path string) (ServerConfig, error) {
    var serverConfig ServerConfig

    config, err := readConfig(path)
    if err != nil {
        return serverConfig, err
    }

    serverConfig.Port = config["Port"].(string)

    privateKey, err := parsePrivateKey(config["PrivateKey"].(string))
    if err != nil {
        panic(err)
    }
    serverConfig.PrivateKey = privateKey

    publicKey, err := parsePublicKey(config["PublicKey"].(string))
    if err != nil {
        panic(err)
    }
    serverConfig.PublicKey = publicKey

    serverConfig.DatabasePath = config["DatabasePath"].(string)

    return serverConfig, nil
}

func readConfig(path string) (map[string]any, error) {
    config := make(map[string]any)

    file, err := os.ReadFile(path)
    if err != nil {
        return config, err
    }

    err = json.Unmarshal(file, &config)

    return config, err
}

func parsePrivateKey(path string) (*rsa.PrivateKey, error) {
    PEM, err := os.ReadFile(path)
    if err != nil {
        return &rsa.PrivateKey{}, err
    }

    block, _ := pem.Decode([]byte(PEM))
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

    if err != nil {
        fmt.Println("Happy-files currently accepts only RSA keys (PKCS1)")
    }

    return privateKey, err
}

func parsePublicKey(path string) (*rsa.PublicKey, error) {
    PEM, err := os.ReadFile(path)
    if err != nil {
        return &rsa.PublicKey{}, err
    }

    block, _ := pem.Decode([]byte(PEM))
    publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)

    if err != nil {
        fmt.Println("Happy-files currently accepts only RSA keys (PKCS1)")
    }

    return publicKey, err
}

func MarshalPrivateKey(key *rsa.PrivateKey) string {
    pemKeyBytes := x509.MarshalPKCS1PrivateKey(key)
    pemKey := pem.EncodeToMemory(&pem.Block{
        Type: "RSA PRIVATE KEY",
        Bytes: pemKeyBytes,
    })
    
    return string(pemKey)
}

func MarshalPublicKey(key *rsa.PublicKey) string {
    pemKeyBytes := x509.MarshalPKCS1PublicKey(key)
    pemKey := pem.EncodeToMemory(&pem.Block{
        Type: "RSA PUBLIC KEY",
        Bytes: pemKeyBytes,
    })
    
    return string(pemKey)
}

func GenKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
    if err != nil {
        panic(err)
    }

    return privateKey, &privateKey.PublicKey, nil
}
