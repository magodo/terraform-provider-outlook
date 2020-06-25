package msauth_test

import (
	"context"
	"testing"

	"github.com/magodo/terraform-provider-outlook/msauth"
)

func TestObtainTokenViaClientCredentials(t *testing.T) {
	clientID := "6731de76-14a6-49ae-97bc-6eba6914391e" // msgraph tutorial client id
	clientCredential := `JqQX2PNo9bpM0uEihUPzyrh`
	c := msauth.NewClientViaClientCredential("common", clientID, clientCredential, "https://graph.microsoft.com/.default")
	_, err := c.ObtainTokenSource(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
