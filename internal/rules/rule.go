package rules

import "github.com/Kukuruzoed/config-analyzer/internal/config"

type Severity string

type Config = config.Config

const (
	LOW    Severity = "LOW"
	MEDIUM Severity = "MEDIUM"
	HIGH   Severity = "HIGH"
)

type Issue struct {
	Severity       Severity
	Description    string
	Recommendation string
}

type Rule interface {
	Check(Config) []Issue
}
