package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/viper"
	"goji.io"
	"goji.io/pat"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "domainName")
	ip := pat.Param(r, "ipAddress")
	fmt.Fprintf(w, "Updating %s to %s", name, ip)
}

func main() {
	err := initConfig()
	if err != nil {
		fmt.Println("Init error:", err)
		return
	}

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/drmu/update/:domainName/:ipAddress"), updateHandler)

	http.ListenAndServe("localhost:8000", mux)
}

func initConfig() error {
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
