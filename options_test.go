package whoisapi

import (
	"net/url"
	"reflect"
	"testing"
)

//TestOptions tests the Options functions
func TestOptions(t *testing.T) {

	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "output format",
			values: url.Values{},
			option: OptionOutputFormat("JSON"),
			want:   "outputFormat=JSON",
		},
		{
			name:   "prefer fresh",
			values: url.Values{},
			option: OptionPreferFresh(1),
			want:   "preferFresh=1",
		},
		{
			name:   "domain availability",
			values: url.Values{},
			option: OptionDA(2),
			want:   "da=2",
		},
		{
			name:   "IP",
			values: url.Values{},
			option: OptionIP(1),
			want:   "ip=1",
		},
		{
			name:   "IP Whois",
			values: url.Values{},
			option: OptionIPWhois(0),
			want:   "ipWhois=0",
		},
		{
			name:   "check proxy data",
			values: url.Values{},
			option: OptionCheckProxyData(1),
			want:   "checkProxyData=1",
		},
		{
			name:   "thin whois",
			values: url.Values{},
			option: OptionThinWhois(0),
			want:   "thinWhois=0",
		},
		{
			name:   "ignore raw texts",
			values: url.Values{},
			option: OptionIgnoreRawTexts(1),
			want:   "ignoreRawTexts=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
