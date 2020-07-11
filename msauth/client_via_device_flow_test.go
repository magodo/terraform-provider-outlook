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

	scopes := []string{
		"mailboxsettings.readwrite",
		"mail.readwrite",
		"offline_access",
	}

	clientID := "23bd8cd9-a50b-4839-b522-67b77d5db7da"
	c := msauth.NewClientViaDeviceFlow("common", clientID, nil, scopes...)
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
