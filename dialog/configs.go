package dialog

import (
	"encoding/json"
	"fmt"
	"os"
)

var defaultClientConfig = map[string]any{
    "username": "",
    "email": "",
    "private-key": "",
    "server-ip": "localhost",
    "server-port": "13333",
    "server-public": "",
    "download-path": "~/Downloads",
}

var defaultServerConfig = map[string]any{
    "private-key": "",
    "db-path": "./hf.db",
}

func readConfig(configPath string) (map[string]any, error) {
    var config map[string]any

    data, err := os.ReadFile(configPath)
    if err != nil {
        return config, err
    }

    err = json.Unmarshal(data, &config)
    if err != nil {
        return config, err
    }

    return config, nil
}

func WriteConfig(configPath string, changes map[string]any) error {
    config, err := readConfig(configPath)
    if err != nil {
        return err
    }

    for k, v := range changes {
        config[k] = v
    }

    
    data, err := json.MarshalIndent(config, "", "\t")
    if err != nil {
        return err
    }

    err = os.WriteFile(configPath, data, 0644)
    if err != nil {
        return err
    }

    return nil
}

func PrintConfig(configPath string) {
    config, err := readConfig(configPath)
    if err != nil {
        panic(err)
    }

    for k, v := range config {
        fmt.Println(k + ": " + v.(string))
    }
    fmt.Printf("\n")
}
