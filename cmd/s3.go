package cmd

import (
	"context"
	"fmt"

	"github.com/Mohamed-M-Meth/aws-dr-audit/internal/auditors"
	awsclients "github.com/Mohamed-M-Meth/aws-dr-audit/internal/aws"
	"github.com/Mohamed-M-Meth/aws-dr-audit/internal/reporter"
	"github.com/spf13/cobra"
)

var s3Cmd = &cobra.Command{
	Use:                    "s3",
	Aliases:                nil,
	SuggestFor:             nil,
	Short:                  "Audit S3 buckets for versioning and Cross-Region Replication",
	GroupID:                "",
	Long:                   "",
	Example:                "",
	ValidArgs:              nil,
	ValidArgsFunction:      nil,
	Args:                   nil,
	ArgAliases:             nil,
	BashCompletionFunction: "",
	Deprecated:             "",
	Annotations:            nil,
	Version:                "",
	PersistentPreRun:       nil,
	PersistentPreRunE:      nil,
	PreRun:                 nil,
	PreRunE:                nil,
	Run:                    nil,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Auditing S3 | Primary: %s Fallback: %s\n\n", region, fallbackRegion)
		ctx := context.Background()
		clients, err := awsclients.NewClients(ctx, region, fallbackRegion, profile)
		if err != nil {
			return fmt.Errorf("failed to initialize AWS clients: %w", err)
		}
		results, err := auditors.AuditS3(ctx, clients, fallbackRegion)
		if err != nil {
			return fmt.Errorf("s3 audit failed: %w", err)
		}
		return reporter.Render(results, outputFormat)
	},
	PostRun:                    nil,
	PostRunE:                   nil,
	PersistentPostRun:          nil,
	PersistentPostRunE:         nil,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	CompletionOptions:          cobra.CompletionOptions{},
	TraverseChildren:           false,
	Hidden:                     false,
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
}
