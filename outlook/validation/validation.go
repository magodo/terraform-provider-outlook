package validation

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		v := i.(string) // this is guaranteed by the schema to be a string

		for _, str := range valid {
			if v == str || (ignoreCase && strings.ToLower(v) == strings.ToLower(str)) {
				return nil
			}
		}

		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected %s to be one of %v, got %s", path, valid, v),
				AttributePath: path,
			},
		}
	}
}
