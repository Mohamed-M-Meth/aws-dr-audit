package cmd

import (
	. "fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	region         string
	fallbackRegion string
	outputFormat   string
	profile        string
)

var rootCmd = &cobra.Command{
	Use:   "aws-dr-audit",
	Short: "AWS Disaster Recovery Audit Tool",
	Long: color.CyanString(`
╔══════════════════════════════════════════════════════════╗
║           AWS-DR-AUDIT — Resilience Scanner              ║
║  Evaluate your AWS account against regional DR failures  ║
╚══════════════════════════════════════════════════════════╝
`) + `
aws-dr-audit scans your AWS infrastructure and reports
Disaster Recovery gaps that would cause data loss during
a complete regional failure (e.g., me-central-1 outage).
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "me-central-1", "Primary AWS region to audit")
	rootCmd.PersistentFlags().StringVarP(&fallbackRegion, "fallback", "f", "eu-central-1", "Target DR fallback region")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table | json")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "AWS credentials profile")

	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(s3Cmd)
	rootCmd.AddCommand(rdsCmd)
}
