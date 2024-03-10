package dialog

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
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

    return client, nil
}

func isAlive(conn net.Conn) bool {
    if conn == nil {
        return false
    } else if _, err := conn.Write(genMsg(PING, PING_FIELD)); errors.Is(err, net.ErrClosed) || errors.Is(err, syscall.EPIPE) {
        return false
    }

    return true
}

func isListening(listener net.Listener) bool {
    if listener == nil {
        return false
    }

    // this probably isn't the right way to do it - what if listener is created but not accepting?
    if _, err := listener.Accept(); err != nil {
        if OpErr, ok := err.(*net.OpError); ok && OpErr.Err.Error() == "use of closed network connection" {
            return false
        } else {
            fmt.Println("Error in isListening(): " + err.Error())
        }
    }

    return true
}


func (client *Client) connectWithServer() error {
    var err error
    client.conn, err = net.Dial(SERVER_CONN_TYPE, client.serverIP + ":" + client.serverPort)
    if err != nil {
        return err
    }
    fmt.Println("Connected to server", client.serverIP + ":" + client.serverPort)

    return nil
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

func (client *Client) download(conn net.Conn, initMsg []byte) error {
    // handles incoming packets: reads, decrypts and appends to file
    defer conn.Close()
    msgFields := groupMsg(initMsg, len(initMsg))
    fmt.Println("Downloading file " + string(msgFields[0]) + "...")
    return nil
}

func (client *Client) updateServer() error {
    if !isAlive(client.conn) {
        if err := client.connectWithServer(); err != nil {
            return err
        }
    }

    listeningAddr := client.listener.Addr().String()
    msg := genMsg(UPDATE_IP, client.username, listeningAddr)
    _, err := client.conn.Write(msg)
    
    if err != nil {
        return err
    }

    return client.conn.Close()
}

func (client *Client) SignUp() error {
    if !isAlive(client.conn) {
        if err := client.connectWithServer(); err != nil {
            return err
        }
    }

    msg := genMsg(SIGN_UP, client.username, client.email)
    if _, err := client.conn.Write(msg); err != nil {
        return err
    }

    return client.conn.Close()
}

func (client *Client) Listen() error {
    var err error
    client.listener, err = net.Listen(TRANSFER_CONN_TYPE, ":0")
    if err != nil {
        return err
    }

    fmt.Println("Listening on:", client.listener.Addr().String())

    if err = client.updateServer(); err != nil {
        return err
    }

    for {
        conn, err := client.listener.Accept()

        if err != nil {
            return err
        }

        // TODO: verify identity

        msgData := make([]byte, 3)
        if _, err := conn.Read(msgData); err != nil {
            return err
        }

        msgCode := msgData[0]
        initMsg := make([]byte, rLength(msgData[1:]))

        if _, err := conn.Read(initMsg); err != nil {
            return err
        }

        switch msgCode {
        case TRANSFER_REQUEST:
            go client.download(conn, initMsg)

        default:
            fmt.Println("Unknown message code:", msgCode)
            conn.Close()
        }
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
