package dialog

import (
	"crypto/rsa"
	"fmt"
	"net"

    // "github.com/happy-files/auth"
)

type Client struct {
    Username, Nickname string
    private_key rsa.PrivateKey
    dest_public_key rsa.PublicKey
    socket net.Conn
}

func (c *Client) NewDialog(dest_ip, dest_port string) {
    fmt.Println("Connecting to", dest_ip + ":" + dest_port, "...")
    connection, err := net.Dial(CONN_TYPE, dest_ip + ":" + dest_port)

    if err != nil {
        panic(err)
    }

    fmt.Println("Connection established")
    c.socket = connection
}

func (c *Client) SignUp(passwd string) error {
    // send (username, hpasswd, public key, nickname)
    // hpasswd := auth.HashPasswd(passwd)
    test := []byte{1}
    test = append(test, ("|" + c.Username)...)
    test = append(test, ("|" + c.Nickname)...)
    test = append(test, ("|" + passwd)...)
    fmt.Println(test)

    c.socket.Write(test)
    return nil
}

func (c *Client) LogIn(passwd string) error {
    // get friends
    // server will update (IP, port)
    packet := []byte{LOG_IN}
    packet = append(packet, passwd...)
    c.socket.Write(packet)
    return nil
}

func (c *Client) getFriendData(friend_username string) (string, string, error) {
    // ask server about friend's public key and (IP, port)
    return "", "", nil
}

func (c *Client) SendFile(friend_username string, file_path string) error {
    return nil
}

func (c *Client) GetFile() error {
    // connection listener
    return nil
}

func (c *Client) Close() error {
    return c.socket.Close()
}
