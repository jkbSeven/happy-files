package main

import (
	"github.com/jkbSeven/happy-files/dialog"
)

func main() {
    client, err := dialog.NewClient("./config.json")
    if err != nil {
        panic(err)
    }
    client.SignUp()
    client.Listen()
}
