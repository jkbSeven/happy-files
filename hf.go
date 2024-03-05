package main

import (
	"github.com/happy-files/dialog"
)

func main() {
    client, err := dialog.NewClient("./config.json")
    if err != nil {
        panic(err)
    }
    client.PrintConfig("./config.json")
    client.Listen()
}
