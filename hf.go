package main

import (
	"github.com/happy-files/dialog"
)

func main() {
    var client dialog.Client

    client.Username = "Username"
    client.Nickname = "Nickname"

    client.NewDialog("localhost", "13333")
    client.SignUp("passwd")
}
