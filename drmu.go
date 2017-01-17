package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/viper"
	"goji.io"
	"goji.io/pat"
)

var cfg = viper.New()

func updateHandler(w http.ResponseWriter, r *http.Request) {
	domain_name := pat.Param(r, "domainName")
	ip := pat.Param(r, "ipAddress")
	client := getClient()
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
		HostedZoneId: aws.String(cfg.GetString("hostedzone")),
	}
	resp, err := client.ChangeResourceRecordSets(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}

func main() {
	mux, err := initConfig()

	if err != nil || mux == nil {
		fmt.Println("Failed to initialize configuration:", err)
		return
	}

	http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.GetString("listenaddress"), cfg.GetString("listenport")), mux)
}

func initConfig() (*goji.Mux, error) {
	cfg.SetConfigName("app")
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("config")
	err := cfg.ReadInConfig()

	if err != nil {
		fmt.Println("Failed to get config:", err)
		return nil, err
	}

	mux := goji.NewMux()

	if mux == nil {
		fmt.Println("Failed to instantiate HTTP routing engine:", err)
		return nil, err
	} else {
		mux.HandleFunc(pat.Get("/drmu/update/:domainName/:ipAddress"), updateHandler)
	}

	return mux, nil
}

func getClient() *route53.Route53 {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Failed to create session:", err)
		return nil
	}

	client := route53.New(sess)

	if err != nil {
		fmt.Println("Failed to create client:", err)
		return nil
	}

	return client
}
