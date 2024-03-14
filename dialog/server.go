package dialog

import (
	"crypto/rsa"
	"fmt"
	"net"
)

type Server struct {
    private_key rsa.PrivateKey
    users map[string]string
}

func (server *Server) Run(port string) {
    server.users = make(map[string]string)
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

        msgCode, response, err := readMsg(conn)
        if err != nil {
            panic(err)
        }

        switch msgCode {
        case SIGN_UP:
            go server.signUp(conn, response)
        case PING:
            go respondPing(conn)
        case UPDATE_IP:
            go server.updateClient(conn, response)
        case GET_USER_DATA:
            go server.userData(conn, response)
        default:
            fmt.Printf("Unknown message code: %d\n", msgCode)
        }

    }
}

func (server *Server) signUp(conn net.Conn, msg []byte) error {
    defer conn.Close()

    fmt.Println("Proccessing SignUp for", conn.RemoteAddr().String()) 

    msgFields := groupMsg(msg, len(msg))
    fmt.Printf("username: " + string(msgFields[0]) + "\nemail: " + string(msgFields[1]) + "\n")

    // TODO: add to database
    // err := sql.query(...)
    // if err != nil ...

    return nil
}

func (server *Server) updateClient(conn net.Conn, msg []byte) error {
    msgFields := groupMsg(msg, len(msg))
    username, addr := string(msgFields[0]), string(msgFields[1])
    server.users[username] = addr

    return nil
}

func (server *Server) userData(conn net.Conn, msg []byte) error {
    msgFields := groupMsg(msg, len(msg))
    dest := string(msgFields[1])

    destIP, ok := server.users[dest]
    if !ok {
        response := genMsg(ERROR, "User does not exist")
        conn.Write(response)
        return conn.Close()
    }

    response := genMsg(GET_USER_DATA, destIP)
    if _, err := conn.Write(response); err != nil {
        return err
    }
    return conn.Close()
}
