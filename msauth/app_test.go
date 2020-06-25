// +build ignore

package msauth_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/magodo/terraform-provider-outlook/msauth"
)

func TestAppObtainTokenViaClientCredentials(t *testing.T) {
	if os.Getenv(EnvvarTenantID) == "" || os.Getenv(EnvvarClientID) == "" || os.Getenv(EnvvarClientCredential) == "" {
		t.Skip(fmt.Sprintf(
			"Test skipped unless env '%s', '%s' and '%s' are set",
			EnvvarTenantID, EnvvarClientID, EnvvarClientCredential))
		return
	}

	app, err := msauth.NewApp(msauth.NewClientViaClientCredential(
		msauth.NewHTTPClient(
			http.DefaultClient, nil),
		msauth.NewScope("https://graph.microsoft.com/.default"),
		os.Getenv(EnvvarClientID),
		os.Getenv(EnvvarClientCredential),
	),
		fmt.Sprintf("https://login.microsoftonline.com/%s", os.Getenv(EnvvarTenantID)))
	if err != nil {
		t.Fatal(err)
	}
	_, err = app.ObtainToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	_, err = app.ObtainTokenSilently(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestAppObtainTokenViaDeviceCode(t *testing.T) {
	if os.Getenv(EnvvarInteractive) == "" {
		t.Skip(fmt.Sprintf(
			"Test skipped unless env '%s' is set",
			EnvvarInteractive))
		return
	}

	// User can either use first party client (as shwon below) to authorize against "common" authority,
	// or use a self registerd client to authorize against "<tenant>" authority.
	clientID := "6731de76-14a6-49ae-97bc-6eba6914391e" // msgraph tutorial client id
	app, err := msauth.NewApp(
		msauth.NewClientViaDeviceFlow(
			msauth.NewHTTPClient(http.DefaultClient, nil),
			msauth.NewScope("user.read", "offline_access"),
			clientID,
			nil),
		fmt.Sprintf("https://login.microsoftonline.com/%s", os.Getenv(EnvvarTenantID)))
	if err != nil {
		t.Fatal(err)
	}
	_, err = app.ObtainToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	_, err = app.ObtainTokenSilently(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
