package config

import (
	"testing"
)


const (
    username = "user"
    email = "user@hf.go"
    storehist = true
    downloadPath = "~/Downloads"
    serverIP = "127.0.0.1:37377"
)

func TestReadClientConfig(t *testing.T) {
    clientConfig, err := ReadClientConfig("example_client_config.json")
    if err != nil {
        t.Fatalf("%v", err)
    }

    if clientConfig.Username != username {
        t.Fatalf("Expected Username to be '%s' and got '%s'", username, clientConfig.Username)
    }
    if clientConfig.Email != email {
        t.Fatalf("Expected Email to be '%s' and got '%s'", email, clientConfig.Email)
    }
    if clientConfig.StoreHistory != storehist {
        t.Fatalf("Expected StoreHistory to be '%t' and got '%t'", storehist, clientConfig.StoreHistory)
    }
    if clientConfig.DownloadPath != downloadPath {
        t.Fatalf("Expected DownloadPath to be '%s' and got '%s'", downloadPath, clientConfig.DownloadPath)
    }
    if clientConfig.ServerIP != serverIP {
        t.Fatalf("Expected ServerIP to be '%s' and got '%s'", serverIP, clientConfig.ServerIP)
    }
}
