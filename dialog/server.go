package dialog

import (
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"strings"
)

type Server struct {
    private_key rsa.PrivateKey
}

func (s *Server) NewDialog(port string) {
    listener, err := net.Listen(SERVER_CONN_TYPE, ":" + port)

    if err != nil {
        panic(err)
    }

    defer listener.Close()

    fmt.Println("Server is running on port", port)
    for {
        conn, err := listener.Accept()

        if err != nil {
            fmt.Println("Error while accepting new connection:", err.Error())
            os.Exit(1)
        } else {
            fmt.Println("Accepted connection from:", conn.RemoteAddr().String())
        }

        msgCode := make([]byte, 1)
        _, err = conn.Read(msgCode)

        if err != nil {
            panic(err)
        }

        switch msgCode[0] {
        case 1:
            go s.signUp(conn)
        default:
            fmt.Printf("Unknown message code: %d\n", msgCode[0])
        }

    }
}

func (s *Server) signUp(conn net.Conn) error {
    fmt.Println("Proccessing SignUp for", conn.RemoteAddr().String()) 
    msg := make([]byte, 2048)

    _, err := conn.Read(msg)
    if err != nil {
        conn.Write(msgBytes(ERROR, err.Error()))
        return err
    }

    userFields := strings.Split(string(msg[1:]), "|")
    username, email, publicKeyE, publicKeyN := userFields[0], userFields[1], userFields[2], userFields[3]
    conn.Write(msgBytes(SIGN_UP, username, email, publicKeyE, publicKeyN))
    
    // TODO: add to database
    // err := sql.query(...)
    // if err != nil ...


    err = conn.Close()
    if err != nil {
        conn.Write(msgBytes(ERROR, err.Error()))
        return err
    }
    return nil
}

