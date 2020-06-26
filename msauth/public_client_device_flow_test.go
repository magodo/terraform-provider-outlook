package msauth_test

import (
	"context"
	"fmt"
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

	clientID := "6731de76-14a6-49ae-97bc-6eba6914391e" // msgraph tutorial client id
	c := msauth.NewPublicClientViaDeviceFlow("common", clientID, nil, "user.read", "offline_access")
	tk, err := c.ObtainToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ts, err := c.ObtainTokenSource(context.Background(), tk)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ts.Token()
	if err != nil {
		t.Fatal(err)
	}
}
