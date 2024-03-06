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

func WriteDefaultClientConfig(configPath string) error {
    config := defaultClientConfig

    data, err := json.Marshal(config)
    if err != nil {
        return err
    }

    return os.WriteFile(configPath, data, 0644)
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

func (client *Client) AddUserToList(listType int, usernames... string) error {
    return nil
}

func (client *Client) UsersFromList(listType int) ([]string, error) {
    return []string{}, nil
}

func (client *Client) userData(username string) ([]string, error) {
    // get (IP, port) and public key of another user
    return []string{}, nil
}

func (client *Client) download(conn net.Conn) error {
    // handles incoming packets: reads, decrypts and appends to file
    defer conn.Close()
    return nil
}

func (client *Client) updateIP() error {
    // makes sure that server has the most recent (IP, port)
    listeningAddr := client.listener.Addr().String()
    msg := listeningAddr
    fmt.Println(msg)
    return nil
}

func (client *Client) SignUp() error {
    msg := genMsg(SIGN_UP, client.username, client.email)
    sent, err := client.conn.Write(msg)
    if err != nil {
        panic(err)
    }
    
    if sent != len(msg) {
        log.Fatalf("Sent %d out of %d bytes during signup", sent, len(msg))
    }

    fmt.Println(msg)

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
        fmt.Println(read)

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
