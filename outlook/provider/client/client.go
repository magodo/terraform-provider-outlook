package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Client interface {
	ListUser(ctx context.Context) ([]User, error)
}

func BuildClient(token string) Client {
	return &client{
		token: token,
	}
}

type client struct {
	token string
}

func (c *client) ListUser(ctx context.Context) ([]User, error) {
	url := "https://graph.microsoft.com/v1.0/users"

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + c.token

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "sending request")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "reading response")
		return nil, err
	}

	var m map[string]*json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		err = errors.Wrap(err, "unmarshalling response to raw message")
		return nil, err
	}
	if m["value"] == nil {
		return nil, errors.New("No `value` in response")
	}

	var users []User
	if err := json.Unmarshal(*m["value"], &users); err != nil {
		err = errors.Wrap(err, "unmarshalling response")
		return nil, err
	}
	return users, nil
}
