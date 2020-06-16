package msauth

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewAuthority(t *testing.T) {
	cases := []struct {
		authURL         string
		expectAuthority *authority
	}{
		{
			"https://login.microsoftonline.com/a20e83fc-34d6-4c8e-8ae7-bf3d5eac71aa",
			&authority{
				AuthorizationEndpoint: "https://login.microsoftonline.com/a20e83fc-34d6-4c8e-8ae7-bf3d5eac71aa/oauth2/v2.0/authorize",
				TokenEndpoint:         "https://login.microsoftonline.com/a20e83fc-34d6-4c8e-8ae7-bf3d5eac71aa/oauth2/v2.0/token",
				Tenant:                "a20e83fc-34d6-4c8e-8ae7-bf3d5eac71aa",
			},
		},
	}

	for _, c := range cases {
		out, err := NewAuthority(c.authURL, http.DefaultClient)
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
