[![whois-api-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![whois-api-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/whois-api-go)
[![whois-api-go test](https://github.com/whois-api-llc/whois-api-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/whois-api-go/actions/)

# Overview

The client library for
[Whois API](https://whois.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/whois-api-go
```

# Examples

Full API documentation available [here](https://whois.whoisxmlapi.com/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := whoisxmlapigo.NewBasicClient(apikey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := whoisxmlapigo.NewClient(apiKey, whoisxmlapigo.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Whois API provides the registration details of a domain name. 

```go

// Make request to get WHOIS record for domain
rec, _, err := client.WhoisService.Data(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

log.Println(rec.DomainName, rec.Audit.UpdatedDate)

// Make request to get raw WHOIS API data
resp, err := client.WhoisService.RawData(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))


```
