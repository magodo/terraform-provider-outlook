package msauth_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/magodo/terraform-provider-outlook/msauth"
)

func TestAppObtainToken(t *testing.T) {
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
}
