package whoisapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathWhoisResponseOK    = "/whois/ok"
	pathWhoisResponseError = "/whois/error"
	// pathWhoisResponseOKwError event can happen, for example, when a domain name is invalid
	pathWhoisResponseOKwError   = "/whois/ok_error"
	pathWhoisResponse500        = "/whois/500"
	pathWhoisResponsePartial1   = "/whois/partial"
	pathWhoisResponsePartial2   = "/whois/partial2"
	pathWhoisResponseUnparsable = "/whois/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// whoisServer is the sample of the Whois API server for testing
func whoisServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathWhoisResponseOK:
		case pathWhoisResponseError:
			w.WriteHeader(400)
			response = respErr
		case pathWhoisResponseOKwError:
			response = respErr
		case pathWhoisResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathWhoisResponsePartial1:
			response = response[:len(response)-10]
		case pathWhoisResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathWhoisResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Whois API client for testing
func newAPI(apiServer *httptest.Server, link string) *Client {

	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}
	apiURL.Path = link

	params := ClientParams{
		HTTPClient:      apiServer.Client(),
		WhoisBaseURL:    apiURL,
		HistoricBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestWhoisAPIData tests the Data function
func TestWhoisAPIData(t *testing.T) {

	checkResultRec := func(res *WhoisRecord) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"WhoisRecord": {
  "createdDate": "2009-03-19T21:47:17Z",
  "updatedDate": "2021-12-26T09:13:06Z",
  "expiresDate": "2027-03-19T21:47:17Z",
  "domainName": "whoisxmlapi.com",
  "status": "clientTransferProhibited clientUpdateProhibited clientRenewProhibited clientDeleteProhibited",
  "parseCode": 3515,
  "audit": {
    "createdDate": "2022-04-07 07:42:54 UTC",
    "updatedDate": "2022-04-07 07:42:54 UTC"
  },
  "registrarName": "GoDaddy.com, LLC",
  "registrarIANAID": "146",
  "contactEmail": "abuse@godaddy.com",
  "domainNameExt": ".com",
  "estimatedDomainAge": 4766
}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><WhoisRecord>
  <createdDate>2009-03-19T21:47:17Z</createdDate>
  <updatedDate>2021-12-26T09:13:06Z</updatedDate>
  <expiresDate>2027-03-19T21:47:17Z</expiresDate>
  <domainName>whoisxmlapi.com</domainName>
</WhoisRecord>`

	const errResp = `{"ErrorMessage": {
  "errorCode": "WHOIS_00",
  "msg": "test error message"
}}`

	server := whoisServer(resp, respUnparsable, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}
	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathWhoisResponseOK,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathWhoisResponse500,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathWhoisResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathWhoisResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathWhoisResponseError,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "API error: [WHOIS_00] test error message",
		},
		{
			name: "unparsable response",
			path: pathWhoisResponseUnparsable,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "successful request with error",
			path: pathWhoisResponseOKwError,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "API error: [WHOIS_00] test error message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api := newAPI(server, tt.path)

			gotRec, _, err := api.Data(tt.args.ctx, tt.args.options)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("Whois.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("Whois.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("Whois.Get() got = %v, expected nil", gotRec)
				}
			}

		})
	}
}

// TestWhoisAPIRawData tests the RawData function
func TestWhoisAPIRawData(t *testing.T) {

	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"WhoisRecord": {
  "createdDate": "2009-03-19T21:47:17Z",
  "updatedDate": "2021-12-26T09:13:06Z",
  "expiresDate": "2027-03-19T21:47:17Z",
  "domainName": "whoisxmlapi.com",
  "status": "clientTransferProhibited clientUpdateProhibited clientRenewProhibited clientDeleteProhibited",
  "parseCode": 3515,
  "audit": {
    "createdDate": "2022-04-07 07:42:54 UTC",
    "updatedDate": "2022-04-07 07:42:54 UTC"
  },
  "registrarName": "GoDaddy.com, LLC",
  "registrarIANAID": "146",
  "contactEmail": "abuse@godaddy.com",
  "domainNameExt": ".com",
  "estimatedDomainAge": 4766
}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><WhoisRecord>
  <createdDate>2009-03-19T21:47:17Z</createdDate>
  <updatedDate>2021-12-26T09:13:06Z</updatedDate>
  <expiresDate>2027-03-19T21:47:17Z</expiresDate>
  <domainName>whoisxmlapi.com</domainName>
</WhoisRecord>`

	const errResp = `{"ErrorMessage": {
  "errorCode": "WHOIS_00",
  "msg": "test error message"
}}`

	server := whoisServer(resp, respUnparsable, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}
	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathWhoisResponseOK,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathWhoisResponse500,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathWhoisResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathWhoisResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathWhoisResponseUnparsable,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathWhoisResponseError,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "API failed with status code: 400",
		},
		{
			name: "successful request with error",
			path: pathWhoisResponseOKwError,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api := newAPI(server, tt.path)

			resp, err := api.RawData(tt.args.ctx, tt.args.options)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("Whois.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !checkResultRaw(resp.Body) {
				t.Errorf("Whois.Get() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
