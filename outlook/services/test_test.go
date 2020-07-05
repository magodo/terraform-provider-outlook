package services_test

import (
	"github.com/magodo/terraform-provider-outlook/msauth"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/terraform-provider-outlook/outlook/provider"
)

func preCheck(t *testing.T) {
	variables := []string{
		"OUTLOOK_TOKEN_CACHE_PATH",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}

	app := msauth.NewApp()
	path := os.Getenv("OUTLOOK_TOKEN_CACHE_PATH")
	if err := app.ImportCache(path); err != nil {
		t.Fatalf("importing auth cache from %s: %+v", path, err)
	}
}

var providerFactories = map[string]func() (*schema.Provider, error){
	"outlook": func() (*schema.Provider, error) {
		return provider.Provider(), nil
	},
}

func importStep(name string, ignore ...string) resource.TestStep {
	step := resource.TestStep{
		ResourceName:      name,
		ImportState:       true,
		ImportStateVerify: true,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}
