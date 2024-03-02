package main

import (
	"github.com/happy-files/dialog"
)

func main() {
    client, err := dialog.NewClient("/dev/null")
    if err != nil {
        panic(err)
    }
    client.SignUp()
}
