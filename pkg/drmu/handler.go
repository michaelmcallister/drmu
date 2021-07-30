package drmu

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/golang/glog"
	"goji.io"
	"goji.io/pat"
)

const (
	recordComment = "Updated by DRMU"
	urlPath       = "/drmu/update/:domainName/:ipAddress"
)

type Client struct {
	route53 *route53.Route53
	mux     *goji.Mux
	zoneID  string
	ttl     int64
}

func New(hostedzone string, route53 *route53.Route53, ttl int64) *Client {
	c := &Client{
		route53: route53,
		zoneID:  hostedzone,
		ttl:     ttl,
		mux:     goji.NewMux(),
	}
	log.V(2).Infof("Adding handler to path %s\n", urlPath)
	c.mux.HandleFunc(pat.Get(urlPath), c.updateHandler)
	return c
}

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.V(2).Infof("New request %s -> %s\n", r.RemoteAddr, r.URL)
	c.mux.ServeHTTP(w, r)
}

func (c *Client) updateHandler(w http.ResponseWriter, r *http.Request) {
	domainName := pat.Param(r, "domainName")
	ip := pat.Param(r, "ipAddress")
	log.V(2).Infof("Creating record %s -> %s\n", domainName, ip)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domainName),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{Value: aws.String(ip)},
						},
						TTL: aws.Int64(c.ttl),
					},
				},
			},
			Comment: aws.String(recordComment),
		},
		HostedZoneId: aws.String(c.zoneID),
	}
	if _, err := c.route53.ChangeResourceRecordSets(params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Unable to update Route53 record: %s", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
