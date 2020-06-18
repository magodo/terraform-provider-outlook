package msauth_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/magodo/terraform-provider-outlook/msauth"
)

func TestObtainTokenViaClientCredentials(t *testing.T) {
	if os.Getenv(EnvvarTenantID) == "" || os.Getenv(EnvvarClientID) == "" || os.Getenv(EnvvarClientCredential) == "" {
		t.Skip(fmt.Sprintf(
			"Test skipped unless env '%s', '%s' and '%s' are set",
			EnvvarTenantID, EnvvarClientID, EnvvarClientCredential))
		return
	}

	client := msauth.NewHTTPClient(http.DefaultClient, nil)
	authority, err := msauth.NewAuthority(fmt.Sprintf("https://login.microsoftonline.com/%s", os.Getenv(EnvvarTenantID)), client.Client)
	if err != nil {
		t.Fatal(err)
	}
	c := msauth.NewClientViaClientCredential(client, msauth.NewScope("https://graph.microsoft.com/.default"), os.Getenv(EnvvarClientID), os.Getenv(EnvvarClientCredential))
	_, err = c.ObtainToken(context.Background(), authority)
	if err != nil {
		t.Fatal(err)
	}
}
