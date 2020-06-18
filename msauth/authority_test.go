package msauth

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
)

func TestNewAuthority(t *testing.T) {
	cases := []struct {
		authURL         string
		expectAuthority *authority
	}{
		{
			"https://login.microsoftonline.com/common",
			&authority{
				AuthorizationEndpoint: "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
				TokenEndpoint:         "https://login.microsoftonline.com/common/oauth2/v2.0/token",
				DeviceEndpoint:        "https://login.microsoftonline.com/common/oauth2/v2.0/devicecode",
			},
		},
	}

	for _, c := range cases {
		out, err := NewAuthority(c.authURL, retryablehttp.NewClient())
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(out, c.expectAuthority) {
			t.Errorf(`Not meet expected authority:
Expect:
%+v
Actual:
%+v
`, c.expectAuthority, out)
		}
	}
}
