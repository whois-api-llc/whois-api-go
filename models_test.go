package whoisapi

import (
	"encoding/json"
	"testing"
)

//TestTime tests JSON encoding/parsing functions for the time values
func TestTime(t *testing.T) {
	tests := []struct {
		name   string
		decErr string
		encErr string
	}{
		{
			name:   `"2006-01-02 15:04:05 EST"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02 12:04:05 UTC"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02T15:04:05-07:00"`,
			decErr: `parsing time "2006-01-02T15:04:05-07:00" as "2006-01-02 15:04:05 MST": cannot parse "T15:04:05-07:00" as " "`,
			encErr: "",
		},
		{
			name:   `""`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v Time

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.name {
				t.Errorf("got = %v, want %v", string(bb), tt.name)
			}
		})
	}
}

//TestTime tests JSON encoding/parsing functions for the Contact struct
func TestContact(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
		decErr string
		encErr string
	}{
		{
			name:   `test-1`,
			input:  `{}`,
			output: `{"name":"","organization":"","street1":"","street2":"","street3":"","street4":"","city":"","state":"","postalCode":"","country":"","countryCode":"","email":"","telephone":"","telephoneExt":"","fax":"","faxExt":"","rawText":"","unparsable":""}`,
			decErr: "",
			encErr: "",
		},
		{
			name: `test-2`,
			input: `{
        "name": "cont-name",
        "organization": "cont-org",
        "street1": "cont-street1",
		"street2": "cont-street2",
		"street3": "cont-street3",
		"street4": "cont-street4",
        "city": "cont-city",
        "state": "cont-state",
        "postalCode": "cont-postalCode",
        "country": "cont-country",
		"countryCode": "cont-countryCode",
        "email": "cont-email",
        "telephone": "cont-telephone",
        "telephoneExt": "cont-telephoneExt",
        "fax": "cont-fax",
        "faxExt": "cont-faxExt",
        "rawText": "cont-rawText",
		"unparsable": "cont-unparsable"
      }`,
			output: `{"name":"cont-name","organization":"cont-org","street1":"cont-street1","street2":"cont-street2","street3":"cont-street3","street4":"cont-street4","city":"cont-city","state":"cont-state","postalCode":"cont-postalCode","country":"cont-country","countryCode":"cont-countryCode","email":"cont-email","telephone":"cont-telephone","telephoneExt":"cont-telephoneExt","fax":"cont-fax","faxExt":"cont-faxExt","rawText":"cont-rawText","unparsable":"cont-unparsable"}`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v Contact

			err := json.Unmarshal([]byte(tt.input), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.output {
				t.Errorf("got  = %v", string(bb))
				t.Errorf("want = %v", tt.output)
			}
		})
	}
}

//TestTime tests JSON encoding/parsing functions for the NameServers struct
func TestNameServers(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
		decErr string
		encErr string
	}{
		{
			name:   `test-1`,
			input:  `{}`,
			output: `{"rawText":"","hostNames":null,"ips":null}`,
			decErr: "",
			encErr: "",
		},
		{
			name: `test-2`,
			input: `{
				"rawText":"",
				"hostNames":["CARL.NS.CLOUDFLARE.COM", "ELLE.NS.CLOUDFLARE.COM"],
				"ips":["104.26.13.210", "172.67.71.123"]
			}`,
			output: `{"rawText":"","hostNames":["CARL.NS.CLOUDFLARE.COM","ELLE.NS.CLOUDFLARE.COM"],"ips":["104.26.13.210","172.67.71.123"]}`,
			decErr: "",
			encErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v NameServers

			err := json.Unmarshal([]byte(tt.input), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.output {
				t.Errorf("got  = %v", string(bb))
				t.Errorf("want = %v", tt.output)
			}
		})
	}
}

// checkErr checks for an error
func checkErr(t *testing.T, err error, want string) {
	if (err != nil || want != "") && (err == nil || err.Error() != want) {
		t.Errorf("error = %v, wantErr %v", err, want)
	}
}
