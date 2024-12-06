package qawaylinter

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
)

// Settings is the root configuration object for the linter.
// The object is divided into individual rule attributes instead of one generic `rules` object.
// This makes JSON deserialization of the configuration easier.
// In addition, it works around limitations in Generics support in Go.
type Settings struct {
	FunctionRules  FunctionRules  `json:"functions"`
	InterfaceRules InterfaceRules `json:"interfaces"`
}

type FunctionRules []FunctionRule[FunctionRuleResults]
type InterfaceRules []InterfaceRule[InterfaceRuleResults]

// Target defines filters for rules.
// Targets allow users to customize to which nodes a rule should apply to.
// For example, interfaces in the domain package may require comments, but interfaces in an internal dev package may not.
type Target struct {
	Packages []string `json:"packages"`
}

// GetMatchingFunctionRules checks if there are any function rules whose target filters match the current node.
func (f FunctionRules) GetMatchingFunctionRules(node ast.Node, pass *analysis.Pass, file *ast.File) FunctionRules {
	var rules FunctionRules
	for _, rule := range f {
		if rule.IsApplicable(node, pass, file) {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetMatchingInterfaceRules checks if there are any interface rules whose target filters match the current node.
func (f InterfaceRules) GetMatchingInterfaceRules(node ast.Node, pass *analysis.Pass, file *ast.File) InterfaceRules {
	var rules InterfaceRules
	for _, rule := range f {
		if rule.IsApplicable(node, pass, file) {
			rules = append(rules, rule)
		}
	}
	return rules
}

// MatchesPackage checks if the given package matches the target.
// Returns true if the full package path starts with any of the target packages.
// For example, if the target is `["example.com/foo"]`, the package `example.com/foo/bar` will match.
func (t Target) MatchesPackage(pkg *types.Package) bool {
	for _, p := range t.Packages {
		if strings.HasPrefix(pkg.Path(), p) {
			return true
		}
	}
	return false
}
