package main

import (
	tfpasses "github.com/bflad/tfproviderlint/passes"
	"github.com/bflad/tfproviderlint/passes/R009"
	tfxpasses "github.com/bflad/tfproviderlint/xpasses"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	var all []*analysis.Analyzer
	all = append(all, tfpasses.AllChecks...)
	all = append(all, tfxpasses.AllChecks...)

	ignored := map[*analysis.Analyzer]interface{}{
		R009.Analyzer: struct{}{},
	}

	var analyzers []*analysis.Analyzer
	for _, analyzer := range all {
		if ignored[analyzer] != nil {
			continue
		}
		analyzers = append(analyzers, analyzer)
	}
	multichecker.Main(analyzers...)
}
