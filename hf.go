package main

import (
    "fmt"
	"github.com/happy-files/dialog"
)

func main() {
    if err := dialog.WriteConfig("./config.json", "username", "user"); err != nil {
        panic(err)
    }
    if err := dialog.WriteConfig("./config.json", "email", "user@hf.go"); err != nil {
        panic(err)
    }
    if err := dialog.WriteConfig("./config.json", "downloadPath", "./"); err != nil {
        panic(err)
    }
    fmt.Println("ok")
    client, err := dialog.NewClient("./config.json")
    if err != nil {
        panic(err)
    }
    client.PrintConfig("./config.json")
}
