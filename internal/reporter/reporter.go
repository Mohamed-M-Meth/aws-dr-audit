package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Mohamed-M-Meth/aws-dr-audit/pkg/models"
	"github.com/fatih/color"
)

var (
	passColor    = color.New(color.FgGreen, color.Bold)
	failColor    = color.New(color.FgRed, color.Bold)
	warnColor    = color.New(color.FgYellow, color.Bold)
	skippedColor = color.New(color.FgCyan)
)

func Render(results []models.AuditResult, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return renderJSON(results)
	default:
		return renderTable(results)
	}
}

func renderTable(results []models.AuditResult) error {
	if len(results) == 0 {
		fmt.Println("No results to display.")
		return nil
	}
	fmt.Println("\n  AWS DISASTER RECOVERY AUDIT REPORT")
	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("%-8s %-25s %-28s %-8s %s\n", "SERVICE", "RESOURCE", "CHECK", "STATUS", "DETAILS")
	fmt.Println(strings.Repeat("-", 90))
	var passed, failed, warned, skipped int
	for _, r := range results {
		fmt.Printf("%-8s %-25s %-28s %-8s %s\n",
			r.Service,
			truncate(r.ResourceID, 24),
			truncate(r.CheckName, 27),
			colorizeStatus(r.Status),
			r.Detail,
		)
		switch r.Status {
		case models.StatusPass:
			passed++
		case models.StatusFail:
			failed++
		case models.StatusWarning:
			warned++
		case models.StatusSkipped:
			skipped++
		}
	}
	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("\nSummary: ")
	_, _ = passColor.Printf(" %d PASSED ", passed)
	_, _ = failColor.Printf(" %d FAILED ", failed)
	_, _ = warnColor.Printf(" %d WARNING ", warned)
	_, _ = skippedColor.Printf(" %d SKIPPED\n\n", skipped)
	if failed > 0 {
		_, _ = failColor.Println("DR POSTURE: CRITICAL")
	} else if warned > 0 {
		_, _ = warnColor.Println("DR POSTURE: AT RISK")
	} else {
		_, _ = passColor.Println("DR POSTURE: HEALTHY")
	}
	return nil
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func colorizeStatus(s models.Status) string {
	switch s {
	case models.StatusPass:
		return passColor.Sprint("PASS")
	case models.StatusFail:
		return failColor.Sprint("FAIL")
	case models.StatusWarning:
		return warnColor.Sprint("WARN")
	case models.StatusSkipped:
		return skippedColor.Sprint("SKIP")
	default:
		return string(s)
	}
}

type jsonReport struct {
	TotalChecks int                  `json:"total_checks"`
	Passed      int                  `json:"passed"`
	Failed      int                  `json:"failed"`
	Warnings    int                  `json:"warnings"`
	Skipped     int                  `json:"skipped"`
	Results     []models.AuditResult `json:"results"`
}

func renderJSON(results []models.AuditResult) error {
	report := jsonReport{TotalChecks: len(results), Results: results}
	for _, r := range results {
		switch r.Status {
		case models.StatusPass:
			report.Passed++
		case models.StatusFail:
			report.Failed++
		case models.StatusWarning:
			report.Warnings++
		case models.StatusSkipped:
			report.Skipped++
		}
	}
	output, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(output))
	return nil
}
