package whoisapi

import (
	"net/url"
	"strconv"
	"strings"
)

// Option adds parameters to the query
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionPreferFresh(0),
	OptionDA(0),
	OptionIP(0),
	OptionIPWhois(0),
	OptionCheckProxyData(0),
	OptionThinWhois(0),
	OptionIgnoreRawTexts(0),
}

// OptionOutputFormat to set Response output format JSON | XML
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionPreferFresh to set parameter for getting the latest WHOIS record even if it's incomplete
func OptionPreferFresh(value int) Option {
	return func(v url.Values) {
		v.Set("preferFresh", strconv.Itoa(value))
	}
}

// OptionDA to set parameter for a quick check on domain availability
func OptionDA(value int) Option {
	return func(v url.Values) {
		v.Set("da", strconv.Itoa(value))
	}
}

// OptionIP to set parameter for returning IPs for the domain name
func OptionIP(value int) Option {
	return func(v url.Values) {
		v.Set("ip", strconv.Itoa(value))
	}
}

// OptionIPWhois to set parameter for returning the WHOIS record for the hosting IP
// if the WHOIS record for the tld of the input domain is not supported
func OptionIPWhois(value int) Option {
	return func(v url.Values) {
		v.Set("ipWhois", strconv.Itoa(value))
	}
}

// OptionCheckProxyData to set parameter for fetching proxy/WHOIS guard data, if it exists
func OptionCheckProxyData(value int) Option {
	return func(v url.Values) {
		v.Set("checkProxyData", strconv.Itoa(value))
	}
}

// OptionThinWhois to set parameter for returning WHOIS data from registry only, without fetching data from registrar
func OptionThinWhois(value int) Option {
	return func(v url.Values) {
		v.Set("thinWhois", strconv.Itoa(value))
	}
}

// OptionIgnoreRawTexts to set parameter for stripping all raw text from the output
func OptionIgnoreRawTexts(value int) Option {
	return func(v url.Values) {
		v.Set("ignoreRawTexts", strconv.Itoa(value))
	}
}
