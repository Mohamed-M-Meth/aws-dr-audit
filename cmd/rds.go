package cmd

import (
	"context"
	"fmt"

	"github.com/Mohamed-M-Meth/aws-dr-audit/internal/auditors"
	awsclients "github.com/Mohamed-M-Meth/aws-dr-audit/internal/aws"
	"github.com/Mohamed-M-Meth/aws-dr-audit/internal/reporter"
	"github.com/spf13/cobra"
)

var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Audit RDS instances for Multi-AZ and cross-region snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🗄️  Auditing RDS | Primary: %s → Fallback: %s\n\n", region, fallbackRegion)

		ctx := context.Background()
		clients, err := awsclients.NewClients(ctx, region, fallbackRegion, profile)
		if err != nil {
			return fmt.Errorf("failed to initialize AWS clients: %w", err)
		}

		results, err := auditors.AuditRDS(ctx, clients, fallbackRegion)
		if err != nil {
			return fmt.Errorf("RDS audit failed: %w", err)
		}

		return reporter.Render(results, outputFormat)
	},
}
