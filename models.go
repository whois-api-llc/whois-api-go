package whoisapi

import (
	"encoding/json"
	"fmt"
	"time"
)

// unmarshalString parses the JSON-encoded data and returns value as a string
func unmarshalString(raw json.RawMessage) (string, error) {
	var val string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// Time is a helper wrapper on time.Time
type Time time.Time

var emptyTime Time

// UnmarshalJSON decodes time as Whois API does
func (t *Time) UnmarshalJSON(b []byte) error {
	str, err := unmarshalString(b)
	if err != nil {
		return err
	}
	if str == "" {
		*t = emptyTime
		return nil
	}
	v, err := time.Parse("2006-01-02 15:04:05 MST", str)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

// MarshalJSON encodes time as Whois API does
func (t Time) MarshalJSON() ([]byte, error) {
	if t == emptyTime {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(t).Format("2006-01-02 15:04:05 MST") + `"`), nil
}

// Audit is part of the Whois API response
// It represents dates when Whois record was added and updated in our database
type Audit struct {
	// CreatedDate is the date this Whois record is collected on whoisxmlapi.com
	CreatedDate Time `json:"createdDate"`

	// UpdatedDate is the date this Whois record is updated on whoisxmlapi.com
	UpdatedDate Time `json:"updatedDate"`
}

// Contact is part of the Whois API response
type Contact struct {
	// Name is the name of the contact
	Name string `json:"name"`

	// Organization is the name of organization
	Organization string `json:"organization"`

	//Street1 is the name of the street
	Street1 string `json:"street1"`

	//Street2 is the name of the street
	Street2 string `json:"street2"`

	//Street3 is the name of the street
	Street3 string `json:"street3"`

	//Street4 is the name of the street
	Street4 string `json:"street4"`

	// City is the name of the city
	City string `json:"city"`

	// State is the name of the city
	State string `json:"state"`

	// PostalCode is a postal code
	PostalCode string `json:"postalCode"`

	// Country is the name of the country
	Country string `json:"country"`

	// PostalCode is a country code
	CountryCode string `json:"countryCode"`

	// Email is an email address
	Email string `json:"email"`

	// Telephone is a phone number
	Telephone string `json:"telephone"`

	// TelephoneExt is the phone extension number
	TelephoneExt string `json:"telephoneExt"`

	// Fax is a fax number
	Fax string `json:"fax"`

	// FaxExt is the fax extension number
	FaxExt string `json:"faxExt"`

	//RawText is the complete raw text of contact's data
	RawText string `json:"rawText"`

	// Unparsable is the part of the raw text that is not parsable by our whois parser
	Unparsable string `json:"unparsable"`
}

// NameServers is part of the Whois API response
type NameServers struct {
	// RawText is the complete raw text of name servers' data
	RawText string `json:"rawText"`

	// HostNames is a list of name servers' hostnames
	HostNames []string `json:"hostNames"`

	// Ips is a list of name servers' IP addresses
	Ips []string `json:"ips"`
}

// RegistryData is part of the Whois API response
type RegistryData struct {
	baseWhoisRecord

	// WhoisServer is the name of Whois server
	WhoisServer string `json:"whoisServer"`

	// ReferralURL is the referral URL
	ReferralURL string `json:"referralURL"`
}

// baseWhoisRecord is the base part of the Whois record
type baseWhoisRecord struct {
	// DomainName is a domain name
	DomainName string `json:"domainName"`

	// CreatedDateNormalized is the normalized form of the date when the domain name was first registered/created
	CreatedDateNormalized Time `json:"createdDateNormalized"`

	// UpdatedDateNormalized is the normalized form of the date when the whois data was updated
	UpdatedDateNormalized Time `json:"updatedDateNormalized"`

	// ExpiresDateNormalized is the normalized form of the date when the domain name will expire
	ExpiresDateNormalized Time `json:"expiresDateNormalized"`

	// CreatedDate is the date when the domain name was first registered/created
	CreatedDate string `json:"createdDate"`

	// UpdatedDate is the date when the whois data was updated
	UpdatedDate string `json:"updatedDate"`

	// ExpiresDate is the date when the domain name will expire
	ExpiresDate string `json:"expiresDate"`

	// Audit is part of the Whois API response
	// It represents dates when Whois record was added and updated in our database
	Audit Audit `json:"audit"`

	//NameServers are name servers or DNS servers for the domain name
	NameServers NameServers `json:"nameServers"`

	// RegistrarName is a registrar name
	// Registrar is an organization or commercial entity that manages the reservation of Internet domain names
	RegistrarName string `json:"registrarName"`

	// RegistrarIANAID is the IANA ID of the registrar
	RegistrarIANAID string `json:"registrarIANAID"`

	// Status is the status code for the domain name
	Status string `json:"status"`

	// RawText is the complete raw text of the whois record
	RawText string `json:"rawText"`

	// ParseCode is a bitmask indicating which fields are parsed in this Whois record
	ParseCode int `json:"parseCode"`

	// Registrant is the owner of the domain name
	// They are the ones who are responsible for keeping the entire Whois contact information up to date
	Registrant Contact `json:"registrant"`

	// AdministrativeContact is the person in charge of the administrative dealings
	// pertaining to the company owning the domain name
	AdministrativeContact Contact `json:"administrativeContact"`

	// TechnicalContact is the person in charge of all technical questions regarding a particular domain name
	TechnicalContact Contact `json:"technicalContact"`

	// BillingContact is the individual who is authorized by the registrant
	// to receive the invoice for domain name registration and domain name renewal fees
	BillingContact Contact `json:"billingContact"`

	// ZoneContact is the person who tends to the technical aspects of maintaining the domain’s name server
	// and resolver software, and database files
	ZoneContact Contact `json:"zoneContact"`

	// Header is the part of the raw text up until the first identifiable field
	Header string `json:"header"`

	// Footer is the part of the raw text after the last identifiable field
	Footer string `json:"footer"`

	// StrippedText includes part of the raw text excluding header and footer
	// this should only include identifiable fields
	StrippedText string `json:"strippedText"`
}

// WhoisRecord is a Whois record
type WhoisRecord struct {
	baseWhoisRecord

	// RegistryData is the Whois record from the domain name registry
	// Each domain name has potentially up to 2 whois record, one from the registry and one from the registrar
	RegistryData RegistryData `json:"registryData"`

	// ContactEmail is the contact email of the Whois record
	ContactEmail string `json:"contactEmail"`

	// DomainAvailability is the result of checking on domain name availability
	DomainAvailability string `json:"domainAvailability"`

	// DomainNameExt is the domain name extension/suffix
	DomainNameExt string `json:"domainNameExt"`

	// EstimatedDomainAge is the estimated age of the domain in days
	EstimatedDomainAge int `json:"estimatedDomainAge"`

	// Ips is a list of IP addresses for a domain name
	Ips []string `json:"ips"`

	// Custom1FieldName is the name of the custom field detected by our parser
	Custom1FieldName string `json:"custom1FieldName"`

	// Custom1FieldValue is the value of the custom field detected by our parser
	Custom1FieldValue string `json:"custom1FieldValue"`

	// Custom2FieldName is the name of the custom field detected by our parser
	Custom2FieldName string `json:"custom2FieldName"`

	// Custom2FieldValue is the value of the custom field detected by our parser
	Custom2FieldValue string `json:"custom2FieldValue"`

	// Custom3FieldName is the name of the custom field detected by our parser
	Custom3FieldName string `json:"custom3FieldName"`

	// Custom3FieldValue is the value of the custom field detected by our parser
	Custom3FieldValue string `json:"custom3FieldValue"`

	// DataError is the data error text
	DataError string `json:"dataError"`

	// SubRecords are sub-records for this Whois record
	SubRecords []WhoisRecord `json:"subRecords"`
}

// ErrorMessage is an error message
type ErrorMessage struct {
	// ErrorCode is the error code
	ErrorCode string `json:"errorCode"`

	// Message is the error message text
	Message string `json:"msg"`
}

// Error returns error message as a string
func (e ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%s] %s", e.ErrorCode, e.Message)
}
