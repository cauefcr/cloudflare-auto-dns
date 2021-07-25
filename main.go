package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"time"

	cloudflare "github.com/cloudflare/cloudflare-go"
	pipes "github.com/ebuchman/go-shell-pipes"
)

var zoneName = "caue.site"

func doTheCheckingAndTheUpdating() {
	// Construct a new API object
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if os.Getenv("CF_ZONE_NAME") != "" {
		zoneName = os.Getenv("CF_ZONE_NAME")
	}
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch the zone ID
	id, err := api.ZoneIDByName(zoneName) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	// Fetch zone details
	zone, err := api.ZoneDetails(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	recs, err := api.DNSRecords(ctx, zone.ID, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}

	out, err := pipes.RunString("ifconfig -a | grep inet6 | grep global")
	if err != nil {
		log.Fatal(err)
	}
	rx := regexp.MustCompilePOSIX("([0-9a-zA-Z]{4}:){7}[0-9a-zA-Z]{4}")
	myIP := rx.FindString(out)

	founds := []cloudflare.DNSRecord{}
	for _, rec := range recs {
		if myIP == rec.Content {
			founds = append(founds, rec)
		}
	}
	if len(founds) == 2 {
		log.Println("all is good with the cosmos")
		return
	}
	for _, rec := range recs {
		if rx.MatchString(rec.Content) {
			rec.Content = myIP
			api.UpdateDNSRecord(ctx, zone.ID, rec.ID, rec)
		}
	}
}

func main() {
	for range time.Tick(10 * time.Minute) {
		doTheCheckingAndTheUpdating()
	}
}
