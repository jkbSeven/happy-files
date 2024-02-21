package dialog

import (
	"crypto/rsa"
	"fmt"
	"net"
	"os"
)

type Server struct {
    private_key rsa.PrivateKey
    listener net.Listener
}

func (s *Server) NewDialog(port string) {
    listener, err := net.Listen(CONN_TYPE, ":" + port)

    if err != nil {
        panic(err)
    }

    s.listener = listener
    defer s.listener.Close()

    for {
        fmt.Println("server: waiting for connection")
        connection, err := s.listener.Accept()

        if err != nil {
            fmt.Println("Error while accepting new connection:", err.Error())
            os.Exit(1)
        }

        code := make([]byte, 1)
        _, err = connection.Read(code)

        if err != nil {
            panic(err)
        }

        switch code[0] {
        case 1:
            go s.SignUp(connection)
            
        case 2:
            go s.LogIn(connection)
        }

    }
}

func (s *Server) SignUp(connection net.Conn) error {
    fmt.Println("server: proccessing SignUp for", connection.RemoteAddr().String()) 
    rest := make([]byte, 1024)
    _, err := connection.Read(rest)
    if err != nil {
        panic(err)
    }
    fmt.Println("server: rest of the message ->", string(rest))
    return nil
}

func (s *Server) LogIn(connection net.Conn) error {
    return nil
}

func (s *Server) SendRecipientData() error {
    return nil
}
