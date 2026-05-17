package analyzer

import (
	"github.com/Kukuruzoed/config-analyzer/internal/config"
	"github.com/Kukuruzoed/config-analyzer/internal/rules"
)

type Rule = rules.Rule
type Issue = rules.Issue
type Config = config.Config

type Analyzer struct {
	rules []Rule
}

func New(rules []Rule) *Analyzer {
	return &Analyzer{rules: rules}
}

func (a *Analyzer) Analyze(config Config) []Issue {
	var issues = make([]Issue, 0)
	for _, rule := range a.rules {
		issues = append(issues, rule.Check(config)...)
	}
	return issues
}

func DefaultRules() []Rule {
	return []Rule{
		rules.DebugModeRule{},
		rules.OpenBindRule{},
		rules.PlainPasswordRule{},
		rules.WeakAlgorithmRule{},
		rules.TLSRule{},
	}
}
