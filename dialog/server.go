package dialog

import (
	"crypto/rsa"
	"fmt"
	"net"
)

type Server struct {
    private_key rsa.PrivateKey
}

func (server *Server) Run(port string) {
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
            go server.signUp(conn, initMsg)
        case UPDATE_IP:
            go server.updateClient(conn, initMsg)
        default:
            fmt.Printf("Unknown message code: %d\n", msgCode)
        }

    }
}

func (server*Server) signUp(conn net.Conn, initMsg []byte) error {
    defer conn.Close()

    fmt.Println("Proccessing SignUp for", conn.RemoteAddr().String()) 

    msgFields := groupMsg(initMsg, len(initMsg))
    fmt.Printf("username: " + string(msgFields[0]) + "\nemail: " + string(msgFields[1]) + "\n")

    // TODO: add to database
    // err := sql.query(...)
    // if err != nil ...

    return nil
}

func (server *Server) updateClient(conn net.Conn, initMsg []byte) error {
    // testing purposes, not the real functionality of this func
    msgFields := groupMsg(initMsg, len(initMsg))
    addr := string(msgFields[1])

    sendconn, err := net.Dial(TRANSFER_CONN_TYPE, addr)
    if err != nil {
        return err
    }
    defer sendconn.Close()
    sendconn.Write(genMsg(TRANSFER_REQUEST, "test.txt"))
    
    return nil
}
