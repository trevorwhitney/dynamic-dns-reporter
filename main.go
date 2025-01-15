package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

func main() {
	accountId := os.Args[1]
	token := os.Args[2]

	var subDomain string
	if len(os.Args) > 3 {
		subDomain = os.Args[3]
	}

	ctx := context.Background()
	tc := dnsimple.StaticTokenHTTPClient(ctx, token)

	updater := &Updater{
		client:    dnsimple.NewClient(tc),
		ctx:       ctx,
		accountId: accountId,
		zoneName:  "trevorwhitney.net",
	}

	updater.Update(subDomain)
}

type Updater struct {
	client    *dnsimple.Client
	ctx       context.Context
	accountId string
	zoneName  string
}

func (u *Updater) Update(subDomain string) {
	zonesResponse, err := u.client.Zones.ListRecords(u.ctx, u.accountId, u.zoneName, nil)
	if err != nil {
		panic(err)
	}

	for _, zone := range zonesResponse.Data {
		if zone.Type == "A" && (zone.Name == subDomain ||
			(subDomain == "" && zone.Name == u.zoneName)) {
			ip := publicIp()
			if zone.Content != ip {
				fmt.Printf("Updating record %d %s %s %s -> %s\n", zone.ID, zone.Type, zone.Name, zone.Content, ip)
				u.updateZone(zone, ip)
			} else {
				fmt.Printf("Record %d %s %s %s is up to date\n", zone.ID, zone.Type, zone.Name, zone.Content)
			}
		}
	}
}

func (u *Updater) updateZone(zone dnsimple.ZoneRecord, ip string) {
	attributes := dnsimple.ZoneRecordAttributes{
		Content: ip,
	}

	_, err := u.client.Zones.UpdateRecord(u.ctx, u.accountId, u.zoneName, zone.ID, attributes)
	if err != nil {
		panic(err)
	}
}

func publicIp() string {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://api.ipify.org", nil)
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return string(resp)
}
