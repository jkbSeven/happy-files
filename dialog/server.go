package dialog

import (
	"crypto/rsa"
	"fmt"
	"net"
)

type Server struct {
    private_key rsa.PrivateKey
}

func (s *Server) Run(port string) {
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
            panic(err)
        } else {
            fmt.Println("Accepted connection from:", conn.RemoteAddr().String())
        }

        msgData := make([]byte, 3)

        if _, err := conn.Read(msgData); err != nil {
            panic(err)
        }

        msgCode := msgData[0]

        initMsg := make([]byte, rLength(msgData[1:]))

        if _, err := conn.Read(initMsg); err != nil {
            panic(err)
        }

        switch msgCode {
        case SIGN_UP:
            go s.signUp(conn, initMsg)
        default:
            fmt.Printf("Unknown message code: %d\n", msgCode)
        }

    }
}

func (s *Server) signUp(conn net.Conn, initMsg []byte) error {
    defer conn.Close()

    fmt.Println("Proccessing SignUp for", conn.RemoteAddr().String()) 

    start, end := 0, FIELD_PREFIX_LEN

    usernameLengthBytes := initMsg[start:end]
    start = end
    end += rLength(usernameLengthBytes)

    username := string(initMsg[start:end])
    start = end
    end += FIELD_PREFIX_LEN

    emailLengthBytes := initMsg[start:end]
    start = end
    end += rLength(emailLengthBytes)

    email := string(initMsg[start:end])


    fmt.Printf("username: " + username + "\nemail: " + email + "\n")

    // TODO: add to database
    // err := sql.query(...)
    // if err != nil ...

    return nil
}

