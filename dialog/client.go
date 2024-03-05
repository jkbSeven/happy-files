package dialog

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Client struct {
    username, email, downloadPath, serverIP, serverPort string
    privateKey rsa.PrivateKey
    serverPublicKey rsa.PublicKey
    conn net.Conn // communicate with server
    listener net.Listener // communicate with other clients
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

func writeClientDefault(configPath string) error {
    config := defaultClient

    data, err := json.Marshal(config)
    if err != nil {
        return err
    }

    return os.WriteFile(configPath, data, 0644)
}

func WriteClient(configPath string, changes map[string]any) error {
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

func NewClient(configPath string) (Client, error) {
    config, err := readConfig(configPath)
    if err != nil {
        panic(err)
    }

    client := Client{
        username: config["username"].(string),
        email: config["email"].(string),
        downloadPath: config["download-path"].(string),
        serverIP: config["server-ip"].(string),
        serverPort: config["server-port"].(string),
    }

    conn, err := net.Dial(SERVER_CONN_TYPE, client.serverIP + ":" + client.serverPort)
    if err != nil {
        panic(err)
    }

    client.conn = conn
    return client, nil
}

func (client *Client) sync() error {
    // sync whitelist and blacklist
    return nil
}

func (client *Client) getUserData(username string) error {
    // get (IP, port) and public key of another user
    return nil
}

func (client *Client) download(conn net.Conn) error {
    // handles incoming packets: reads, decrypts and appends to file
    defer conn.Close()
    return nil
}

func (client *Client) updateServer() error {
    // makes sure that server has the most recent (IP, port)
    return nil
}

func (client *Client) SignUp() error {
    msg := msgBytes(SIGN_UP, client.username, client.email)
    sent, err := client.conn.Write(msg)
    if err != nil {
        panic(err)
    }
    
    if sent != len(msg) {
        log.Fatalf("Sent %d out of %d bytes during signup", sent, len(msg))
    }

    // TODO: verify message signature

    return nil
}

func (client *Client) Listen() error {
    listener, err := net.Listen(TRANSFER_CONN_TYPE, ":0")
    if err != nil {
        panic(err)
    }

    client.listener = listener
    listeningPort := listener.Addr().String()
    fmt.Println("Listening on:", listeningPort)

    // inform server about the port

    for {
        conn, err := listener.Accept()

        if err != nil {
            return err
        }

        // TODO: verify identity

        msg := make([]byte, 3)
        read, err := conn.Read(msg)

        if err != nil {
            return err
        }

        go client.download(conn)
    }
}

func (client *Client) Send(username, filepath string) error {
    // encrypt and send selected file to the user
    return nil
}

func (client *Client) PrintConfig(configPath string) {
    config, err := readConfig(configPath)
    if err != nil {
        panic(err)
    }

    for k, v := range config {
        fmt.Println(k + ": " + v.(string))
    }
    fmt.Printf("\n")
}


func (client *Client) Close() error {
    if err := client.conn.Close(); err != nil {
        return err
    }
    return client.listener.Close()
}
