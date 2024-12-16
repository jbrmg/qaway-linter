package qawaylinter

import (
	"github.com/golangci/plugin-module-register/register"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"
)

// AnalyzerPlugin is the entry point for the linter.
// Please see https://golangci-lint.run/plugins/module-plugins/ for instructions on how to integrate custom linters
// into golangci-lint.
// There is also an example linter at https://github.com/golangci/example-plugin-module-linter which was used
// as baseline for this implementation.
type AnalyzerPlugin struct {
	Settings Settings
}

func (a *AnalyzerPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name:     "qawaylinter",
			Doc:      "Checks that a given function has an appropriate amount of documetation.",
			Run:      a.Run,
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		},
	}, nil
}

// Run executes the analysis step of the linter.
// The method iterates over all files and applies all rules to the nodes in the file.
// Please refer to the Rule interface for more information on how to implement rules.
// Do to limitations in Go generics, the rules are split, e.g. into FunctionRules and InterfaceRules.
// It would be better if they would all be part of a single list as they all implement the same interface,
// but this was not possible.
func (a *AnalyzerPlugin) Run(pass *analysis.Pass) (interface{}, error) {
	var file *ast.File
	inspect := func(node ast.Node) bool {
		functionRules := a.Settings.FunctionRules.GetMatchingFunctionRules(node, pass, file)
		for _, rule := range functionRules {
			results := rule.Analyse(node, pass, file)
			rule.Apply(results, node, pass)
		}

		interfaceRules := a.Settings.InterfaceRules.GetMatchingInterfaceRules(node, pass, file)
		for _, rule := range interfaceRules {
			results := rule.Analyse(node, pass, file)
			rule.Apply(results, node, pass)
		}

		structRules := a.Settings.StructRules.GetMatchingStructRules(node, pass, file)
		for _, rule := range structRules {
			results := rule.Analyse(node, pass, file)
			rule.Apply(results, node, pass)
		}
		return true

	}

	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename

		// skip all tests fails as documenting them is not as important
		if strings.HasSuffix(filename, "_test.go") {
			continue
		}
		file = f
		ast.Inspect(f, inspect)
	}
	return nil, nil
}

func (a *AnalyzerPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
