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

    data, err := json.MarshalIndent(config, "", "\t")
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

func (client *Client) userData(username string) ([][]byte, error) {
    if !isAlive(client.conn) {
        if err := client.connectWithServer(); err != nil {
            return [][]byte{}, err
        }
    }
    
    msg := genMsg(GET_USER_DATA, client.username, username)
    if _, err := client.conn.Write(msg); err != nil {
        return [][]byte{}, err
    }

    msgCode, msg, err := readMsg(client.conn)
    
    if err != nil {
        return [][]byte{}, err
    }

    if msgCode == ERROR {
        return [][]byte{}, errors.New("Error getting user data: " + string(groupMsg(msg, len(msg))[0]))
    } else if msgCode != GET_USER_DATA {
        return [][]byte{}, errors.New("Could not get proper response from the server")
    }

    msgFields := groupMsg(msg, len(msg))
    data := make([][]byte, 3)
    copy(data, msgFields[:3])

    return data, nil
}

func (client *Client) download(conn net.Conn, initMsg []byte) error {
    // handles incoming packets: reads, decrypts and appends to file
    defer conn.Close()
    msgFields := groupMsg(initMsg, len(initMsg))
    fmt.Println("Downloading file " + string(msgFields[1]) + " from " + string(msgFields[0]) + " (size: " + string(msgFields[2]) + " MB)")

    msg := genMsg(TRANSFER_ACCEPTED, PING_FIELD)
    if _, err := conn.Write(msg); err != nil {
        return err
    }

    for {
        msgCode, response, err := readMsg(conn)
        if err != nil {
            return err
        }
        
        if msgCode != TRANSFER {
            return errors.New("Download was interrupted by a packet with a wrong message code")
        }

        msgFields := groupMsg(response, len(response))
        fmt.Println(string(msgFields[0]))
    }
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

        msgCode, response, err := readMsg(conn)
        if err != nil {
            return err
        }

        switch msgCode {
        case TRANSFER_REQUEST:
            go client.download(conn, response)

        default:
            fmt.Println("Unknown message code:", msgCode)
            conn.Close()
        }
    }
}

func (client *Client) Send(username, filepath string) error {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    fileStat, err := file.Stat()
    if err != nil {
        return err
    }

    fileName := fileStat.Name()
    fileSize := fileStat.Size()
    fileSizeMB := int64(float64(fileSize) / 1_000_000)
    fileSizeMBstring := fmt.Sprintf("%d", fileSizeMB)

    userData, err := client.userData(username)
    if err != nil {
        fmt.Println(err)
        return err
    }

    destAddr := string(userData[0])
    // destPubKey := userData[2]

    conn, err := net.Dial(TRANSFER_CONN_TYPE, destAddr)
    if err != nil {
        return err
    }

    msg := genMsg(TRANSFER_REQUEST, client.username, fileName, fileSizeMBstring)
    if _, err := conn.Write(msg); err != nil {
        return err
    }

    msgCode, _, err := readMsg(conn)
    if err != nil {
        return err
    } else if msgCode != TRANSFER_ACCEPTED {
        return errors.New("User denied file transfer")
    }

    for sent := 0; sent < int(fileSize); {
        fileChunk := make([]byte, 512)
        read, err := file.Read(fileChunk)
        if err != nil {
            return err
        }
        
        fileChunkMsg := genMsg(TRANSFER, string(fileChunk[:read]))
        if _, err := conn.Write(fileChunkMsg); err != nil {
            return err
        }

        sent += read
    }

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
    if isAlive(client.conn) {
        if err := client.conn.Close(); err != nil {
            return err
        }
    }
    return client.listener.Close()
}
