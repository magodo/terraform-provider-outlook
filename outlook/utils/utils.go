package utils

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func ResponseErrorWasNotFound(err error) bool {
	if errRes, ok := err.(*msgraph.ErrorResponse); ok {
		return errRes.StatusCode() == http.StatusNotFound
	}
	return false
}

func ImportAsExistsError(resourceName, id string) diag.Diagnostics {
	msg := "A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	return diag.Errorf(msg, id, resourceName)
}
