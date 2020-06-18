package msauth_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/magodo/terraform-provider-outlook/msauth"
)

func TestObtainTokenViaDeviceFlow(t *testing.T) {
	if os.Getenv(EnvvarInteractive) == "" {
		t.Skip(fmt.Sprintf(
			"Test skipped unless env '%s' is set",
			EnvvarInteractive))
		return
	}

	client := msauth.NewHTTPClient(http.DefaultClient, nil)

	// User can either use first party client (as shwon below) to authorize against "common" authority,
	// or use a self registerd client to authorize against "<tenant>" authority.
	clientID := "6731de76-14a6-49ae-97bc-6eba6914391e" // msgraph tutorial client id
	authority, err := msauth.NewAuthority("https://login.microsoftonline.com/common", client.Client)
	if err != nil {
		t.Fatal(err)
	}
	c := msauth.NewClientViaDeviceFlow(client, msauth.NewScope("user.read", "offline_access"), clientID, nil)
	_, err = c.ObtainToken(context.Background(), authority)
	if err != nil {
		t.Fatal(err)
	}
}
