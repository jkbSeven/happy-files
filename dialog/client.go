package dialog

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Client struct {
    username, email, downloadPath string
    privateKey rsa.PrivateKey
    serverPublicKey rsa.PublicKey
    socket net.Conn // communicate with server
    listener net.Listener // communicate with other clients
}

func NewClient(config string) (Client, error) {
    fmt.Println("Reading JSON...")
    conn, err := net.Dial(SERVER_CONN_TYPE, ":13334")
    if err != nil {
        panic(err)
    }
    client := Client{username: "user", email: "user@hf.go", socket: conn}
    return client, nil
}

func lockPort() error {
    return nil
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
    msg := msgBytes(SIGN_UP, client.username, client.email, strconv.Itoa(client.privateKey.PublicKey.E), client.privateKey.N.String())
    fmt.Println("Sent:", string(msg), "-- Length:", len(msg))
    sent, err := client.socket.Write(msg)
    if err != nil {
        panic(err)
    }
    
    if sent != len(msg) {
        log.Fatalf("Sent %d out of %d bytes during signup", sent, len(msg))
    }

    response := make([]byte, 2048)
    _, err = client.socket.Read(response)

    if err != nil {
        panic(err)
    }

    // TODO: verify message signature
    fmt.Println("Received:", string(response), "-- Length:", len(response))

    return nil
}

func (client *Client) Listen(port string) error {
    // start listening on the port that was previously used to communicate with server
    listener, err := net.Listen(TRANSFER_CONN_TYPE, ":0")
    if err != nil {
        panic(err)
    }
    defer listener.Close()

    client.listener = listener

    for {
        conn, err := listener.Accept()

        if err != nil {
            return err
        }

        // TODO: verify identity

        msg := make([]byte, 1024)
        read, err := conn.Read(msg)

        if err != nil {
            return err
        }

        if msg[0] != TRANSFER_REQUEST {
            return errors.New("Message code was different then TRANSFER_REQUEST")
        }

        filename := string(msg[2:read])
        fmt.Println("Downloading " + filename + "...")

        go client.download(conn)
    }
}

func (client *Client) Send(username, filepath string) error {
    // encrypt and send selected file to the user
    return nil
}

func (client *Client) Close() error {
    return client.socket.Close()
}
