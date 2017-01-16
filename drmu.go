package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
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
	client, err := initConfig()
	if err != nil {
		fmt.Println("Failed to initialize configuration:", err)
		return
	}

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/drmu/update/:domainName/:ipAddress"), updateHandler)

	http.ListenAndServe("localhost:8000", mux)
}

func initConfig() (*route53.Route53, error) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

	userConfig := viper.New()
	userConfig.SetConfigName("app")
	userConfig.AddConfigPath(".")
	userConfig.AddConfigPath("config")
	err := userConfig.ReadInConfig()

	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession()

	if err != nil {
		fmt.Println("Failed to create session:", err)
		return nil, err
	}

	client := route53.New(sess)

	if err != nil {
		fmt.Println("Failed to create client:", err)
		return nil, err
	}

	return client, nil
}
