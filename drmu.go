package main

import (
    "fmt"
    "os"

    "github.com/spf13/viper"
)

func main() {
    err := initConfig()
    if err != nil {
        fmt.Println("Init error:", err)
        return
    }
}

func initConfig() (error) {
    os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

    userConfig := viper.New()
    userConfig.SetConfigName("app")
    userConfig.AddConfigPath(".")
    userConfig.AddConfigPath("config")
    err := userConfig.ReadInConfig()
    if err != nil {
        return err
    }
  return nil
}
