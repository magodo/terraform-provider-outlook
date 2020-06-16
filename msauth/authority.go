package msauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type authority struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	Tenant                string
}

func NewAuthority(authorityURL string, client *http.Client) (*authority, error) {
	u, err := url.Parse(authorityURL)
	if err != nil {
		err = errors.Wrapf(err, "parsing authority URL %s", authorityURL)
		return nil, err
	}

	if u.Scheme != "https" {
		return nil, fmt.Errorf("Authority URL %s expect to use HTTPS", authorityURL)
	}

	paths := strings.Split(u.Path, "/")
	if len(paths) != 2 || paths[1] == "" {
		return nil, fmt.Errorf(`Authority URL %s should has the form: "https://<host>/<tenant>"`, authorityURL)
	}
	tenant := paths[1]

	// Discover tenant
	tenantDiscoveryEndpoint := fmt.Sprintf("https://%s%s/v2.0/.well-known/openid-configuration", u.Hostname(), u.Path)
	resp, err := client.Get(tenantDiscoveryEndpoint)
	if err != nil {
		err = errors.Wrapf(err, "GET on %s", tenantDiscoveryEndpoint)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = errors.Wrapf(err, "reading resp from %s", tenantDiscoveryEndpoint)
		return nil, err
	}

	var auth authority

	if err := json.Unmarshal(body, &auth); err != nil {
		return nil, errors.Wrapf(err, "unmarshal %s", body)
	}
	if auth.AuthorizationEndpoint == "" {
		return nil, errors.New(`"authorization_endpoint" is empty`)
	}
	if auth.TokenEndpoint == "" {
		return nil, errors.New(`"token_endpoint" is empty`)
	}
	auth.Tenant = tenant

	return &auth, nil
}
