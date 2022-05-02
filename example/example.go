package example

import (
	"context"
	"errors"
	whoisapi "github.com/whois-api-llc/whois-api-go"
	"log"
	"time"
)

func WhoisData(apikey string) {
	client := whoisapi.NewBasicClient(apikey)

	// Get parsed whois record as a model instance
	rec, resp, err := client.WhoisService.Data(context.Background(), "google.com",
		// this option is ignored, as the inner parser works with JSON only
		whoisapi.OptionOutputFormat("XML"),
		whoisapi.OptionDA(2), whoisapi.OptionIP(1))

	if err != nil {
		// Handle error message returned by server
		var apiErr *whoisapi.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.ErrorCode)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	log.Printf("domainName: %s, audit.updatedDate: %s, updatedDate: %s, domainAvailability: %s, IPs: %s\n",
		rec.DomainName,
		time.Time(rec.Audit.UpdatedDate).Format("2006-01-02 15:04:05 MST"),
		rec.UpdatedDate,
		rec.DomainAvailability,
		rec.Ips)

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func WhoisRawData(apikey string) {
	client := whoisapi.NewBasicClient(apikey)

	// Get raw API response
	resp, err := client.WhoisService.RawData(context.Background(), "whoisxmlapi.com",
		whoisapi.OptionOutputFormat("JSON"), whoisapi.OptionThinWhois(1))

	if err != nil {
		// Handle error message returned by server
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
