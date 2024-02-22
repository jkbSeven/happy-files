package dialog

import (
	"crypto/rsa"
	"fmt"
	"net"
	// "github.com/happy-files/auth"
)

type Client struct {
    username, email string
    private_key rsa.PrivateKey
    socket net.Conn // communicate with server
    listener net.Listener // communicate with other clients
}

func (client *Client) HelloServer(dest_ip, dest_port string) {
    fmt.Println("Connecting to", dest_ip + ":" + dest_port, "...")
    connection, err := net.Dial(TRANSFER_CONN_TYPE, dest_ip + ":" + dest_port)

    if err != nil {
        panic(err)
    }

    fmt.Println("Connection established")
    client.socket = connection
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
    return nil
}

func (client *Client) updateServer() error {
    // makes sure that server has the most recent (IP, port)
    return nil
}

func (client *Client) SignUp() error {
    // sign up with data in config file
    return nil
}

func (client *Client) Listen(port string) error {
    // start listening on the port that was previously used to communicate with server
    return nil
}

func (client *Client) Send(username, filepath string) error {
    // encrypt and send selected file to the user
    return nil
}

func (client *Client) Close() error {
    return client.socket.Close()
}
