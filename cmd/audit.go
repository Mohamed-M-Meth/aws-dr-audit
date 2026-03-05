package cmd

import (
	"context"
	"fmt"

	"github.com/Mohamed-M-Meth/aws-dr-audit/internal/auditors"
	awsclients "github.com/Mohamed-M-Meth/aws-dr-audit/internal/aws"
	"github.com/Mohamed-M-Meth/aws-dr-audit/internal/reporter"
	"github.com/Mohamed-M-Meth/aws-dr-audit/pkg/models"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Run all DR checks (S3 + RDS)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🔍 Starting full DR audit | Primary: %s → Fallback: %s\n\n", region, fallbackRegion)
		ctx := context.Background()
		clients, err := awsclients.NewClients(ctx, region, fallbackRegion, profile)
		if err != nil {
			return fmt.Errorf("failed to initialize AWS clients: %w", err)
		}
		var allResults []models.AuditResult
		s3Results, err := auditors.AuditS3(ctx, clients, fallbackRegion)
		if err != nil {
			return fmt.Errorf("s3 audit failed: %w", err)
		}
		allResults = append(allResults, s3Results...)
		rdsResults, err := auditors.AuditRDS(ctx, clients, fallbackRegion)
		if err != nil {
			return fmt.Errorf("RDS audit failed: %w", err)
		}
		allResults = append(allResults, rdsResults...)
		return reporter.Render(allResults, outputFormat)
	},
}
