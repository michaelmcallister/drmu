package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/viper"
	"goji.io"
	"goji.io/pat"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	domain_name := pat.Param(r, "domainName")
	ip := pat.Param(r, "ipAddress")
	fmt.Fprintf(w, "Updating %s to %s", domain_name, ip)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domain_name),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
						TTL: aws.Int64(300),
					},
				},
			},
			Comment: aws.String("Updated by DRMU"),
		},
		HostedZoneId: aws.String(""),
	}
	resp, err := client.ChangeResourceRecordSets(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}

func main() {
	client, mux, err := initConfig()

	if err != nil || client == nil || mux == nil {
		fmt.Println("Failed to initialize configuration:", err)
		return
	}

	mux.HandleFunc(pat.Get("/drmu/update/:domainName/:ipAddress"), updateHandler)
	http.ListenAndServe("localhost:8000", mux)
}

func initConfig() (*route53.Route53, *goji.Mux, error) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

	userConfig := viper.New()
	userConfig.SetConfigName("app")
	userConfig.AddConfigPath(".")
	userConfig.AddConfigPath("config")
	err := userConfig.ReadInConfig()

	if err != nil {
		return nil, nil, err
	}

	sess, err := session.NewSession()

	if err != nil {
		fmt.Println("Failed to create session:", err)
		return nil, nil, err
	}

	client := route53.New(sess)

	if err != nil {
		fmt.Println("Failed to create client:", err)
		return nil, nil, err
	}

	mux := goji.NewMux()

	if mux == nil {
		fmt.Println("Failed to instantiate HTTP routing engine:", err)
		return nil, nil, err
	} else {
		//mux.HandleFunc(pat.Get("/drmu/update/:domainName/:ipAddress"), updateHandler)
	}

	return client, mux, nil
}
