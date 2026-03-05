package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Clients struct {
	S3Primary      *s3.Client
	RDSPrimary     *rds.Client
	RDSFallback    *rds.Client
	PrimaryRegion  string
	FallbackRegion string
}

func NewClients(ctx context.Context, primaryRegion, fallbackRegion, profile string) (*Clients, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(primaryRegion),
	}
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	primaryCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot load AWS config for region %s: %w", primaryRegion, err)
	}

	fallbackOpts := []func(*config.LoadOptions) error{
		config.WithRegion(fallbackRegion),
	}
	if profile != "" {
		fallbackOpts = append(fallbackOpts, config.WithSharedConfigProfile(profile))
	}
	fallbackCfg, err := config.LoadDefaultConfig(ctx, fallbackOpts...)
	if err != nil {
		return nil, fmt.Errorf("cannot load AWS config for fallback region %s: %w", fallbackRegion, err)
	}

	return &Clients{
		S3Primary:      s3.NewFromConfig(primaryCfg),
		RDSPrimary:     rds.NewFromConfig(primaryCfg),
		RDSFallback:    rds.NewFromConfig(fallbackCfg),
		PrimaryRegion:  primaryRegion,
		FallbackRegion: fallbackRegion,
	}, nil
}
