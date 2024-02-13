package dialog

import (
	"crypto/rsa"
	"net"

	"github.com/happy-files/client/auth"
)

const INIT_CONN_TYPE = "tcp"

type Dialog struct {
    conn net.Conn
    conn_type string
    keep_alive int

    dest_ip, dest_port string
    dest_public_key rsa.PublicKey
}

func NewDialog() Dialog {
    // read data from json, hardcoded for now
    dest_ip, dest_port := "127.0.0.1", "13333"
    connection, err := net.Dial(INIT_CONN_TYPE, dest_ip + ":" + dest_port)

    if err != nil {
        panic(err)
    }

    return Dialog{
        conn: connection,
        conn_type: INIT_CONN_TYPE,
        dest_ip: dest_ip,
        dest_port: dest_port}
}

func (d *Dialog) Introduce(conn net.Conn, username string) error {

    return nil
}

func (d *Dialog) SignUp(username, passwd string, public_key *rsa.PublicKey) error {
    // make connection to the server and send encrypted data (encrypted with server's private key)
    hpasswd := auth.HashPasswd(passwd)
    return nil
}

func (d *Dialog) Login(username, passwd string) error {
    return nil
}
