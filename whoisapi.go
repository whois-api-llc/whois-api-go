package whoisapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// defaultWhoisApiURL is the default Whois API URL
const defaultWhoisApiURL = `https://www.whoisxmlapi.com/whoisserver/WhoisService`

// WhoisService is an interface for Whois API
type WhoisService interface {
	// Data returns parsed Whois record
	Data(ctx context.Context, name string, opts ...Option) (*WhoisRecord, *Response, error)

	// RawData returns raw Whois API response as Response struct with Body saved as a byte slice
	RawData(ctx context.Context, name string, opts ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice
type Response struct {
	*http.Response

	//Body is the byte slice representation of http.Response Body
	Body []byte
}

// whoisApiServiceOp is the type implementing the WhoisService interface
type whoisApiServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ WhoisService = &whoisApiServiceOp{}

// newRequest creates the API request with default parameters and the specified apiKey
func (service *whoisApiServiceOp) newRequest() (*http.Request, error) {

	req, err := service.client.NewRequest(http.MethodGet, service.baseURL, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// whoisApiResponse is used for parsing Whois API response as a model instance
type whoisApiResponse struct {
	WhoisRecord  *WhoisRecord  `json:"WhoisRecord"`
	ErrorMessage *ErrorMessage `json:"ErrorMessage"`
}

// request returns intermediate Whois API response for further actions
func (service *whoisApiServiceOp) request(ctx context.Context, name string, opts ...Option) (*Response, error) {
	if name == "" {
		return nil, &ArgError{"name", "cannot be empty"}
	}

	req, err := service.newRequest()
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("domainName", name)

	for _, opt := range opts {
		opt(q)
	}

	req.URL.RawQuery = q.Encode()

	var b bytes.Buffer
	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw Whois API response
func parse(raw []byte) (*whoisApiResponse, error) {

	var response whoisApiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// Data returns parsed Whois record
func (service whoisApiServiceOp) Data(
	ctx context.Context,
	name string,
	opts ...Option,
) (whois *WhoisRecord, resp *Response, err error) {

	optsJson := make([]Option, 0, len(opts)+1)
	optsJson = append(optsJson, opts...)
	optsJson = append(optsJson, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, name, optsJson...)
	if err != nil {
		return nil, resp, err
	}

	whoisResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if whoisResp.ErrorMessage != nil {
		return nil, nil, ErrorMessage{
			ErrorCode: whoisResp.ErrorMessage.ErrorCode,
			Message:   whoisResp.ErrorMessage.Message,
		}
	}

	return whoisResp.WhoisRecord, resp, nil
}

// RawData returns raw Whois API response as Response struct with Body saved as a byte slice
func (service whoisApiServiceOp) RawData(
	ctx context.Context,
	name string,
	opts ...Option,
) (resp *Response, err error) {

	resp, err = service.request(ctx, name, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
