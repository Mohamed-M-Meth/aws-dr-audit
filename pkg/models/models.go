package models

type Status string

const (
	StatusPass    Status = "PASS"
	StatusFail    Status = "FAIL"
	StatusWarning Status = "WARNING"
	StatusSkipped Status = "SKIPPED"
)

type AuditResult struct {
	Service    string
	ResourceID string
	CheckName  string
	Status     Status
	Detail     string
}
