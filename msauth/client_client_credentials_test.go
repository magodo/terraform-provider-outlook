package msauth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
)

const (
	EnvvarClientID         string = "MSAUTH_CLIENT_ID"
	EnvvarClientCredential string = "MSAUTH_CLIENT_CREDENTIAL"
)

func TestObtainToken(t *testing.T) {
	if os.Getenv(EnvvarClientID) == "" || os.Getenv(EnvvarClientCredential) == "" {
		t.Skip(fmt.Sprintf(
			"Client credentials tests skipped unless env '%s' and '%s' are set",
			EnvvarClientID, EnvvarClientCredential))
		return
	}

	client := NewHTTPClient(http.DefaultClient, nil)
	authority, err := NewAuthority("https://login.microsoftonline.com/a20e83fc-34d6-4c8e-8ae7-bf3d5eac71aa", client.Client)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClientViaClientCredential(client, NewScope("https://graph.microsoft.com/.default"), os.Getenv(EnvvarClientID), os.Getenv(EnvvarClientCredential))
	_, err = c.ObtainToken(context.Background(), authority)
	if err != nil {
		t.Fatal(err)
	}
}
