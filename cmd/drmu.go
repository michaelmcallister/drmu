package main

import (
	"flag"
	"net/http"

	"github.com/michaelmcallister/drmu/pkg/drmu"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/golang/glog"
)

var (
	address      string
	hostedZoneId string
	ttl          int64
)

func init() {
	flag.StringVar(&address, "address", "localhost:8000", "address to bind on")
	flag.StringVar(&hostedZoneId, "zone", "HOSTED_ZONE_ID", "Route53 Hosted Zone ID")
	flag.Int64Var(&ttl, "ttl", 300, "TTL for each record created in Route53")
}

func main() {
	flag.Parse()

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("Unable to create new session: %s", err)
	}
	r53 := route53.New(sess)
	if err != nil {
		log.Fatalf("Unable to create new Route53 client: %s", err)
	}
	c := drmu.New(hostedZoneId, r53, ttl)
	log.Infof("Starting server on %s...", address)
	http.ListenAndServe(address, c)
}
