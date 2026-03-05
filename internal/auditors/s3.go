package auditors

import (
	"context"
	"fmt"
	"strings"

	awsclients "github.com/Mohamed-M-Meth/aws-dr-audit/internal/aws"
	"github.com/Mohamed-M-Meth/aws-dr-audit/pkg/models"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func AuditS3(ctx context.Context, clients *awsclients.Clients, fallbackRegion string) ([]models.AuditResult, error) {
	var results []models.AuditResult
	listOutput, err := clients.S3Primary.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("ListBuckets failed: %w", err)
	}
	if len(listOutput.Buckets) == 0 {
		fmt.Println("  No S3 buckets found.")
		return results, nil
	}
	fmt.Printf("  Found %d S3 bucket(s). Checking each...\n\n", len(listOutput.Buckets))
	for _, bucket := range listOutput.Buckets {
		bucketName := *bucket.Name

		results = append(results, checkS3PublicAccessBlock(ctx, clients.S3Primary, bucketName))

		versioningResult := checkS3Versioning(ctx, clients.S3Primary, bucketName)
		results = append(results, versioningResult)

		if versioningResult.Status == models.StatusPass {
			results = append(results, checkS3CRR(ctx, clients.S3Primary, bucketName, fallbackRegion))
		} else {
			results = append(results, models.AuditResult{
				Service:    "S3",
				ResourceID: bucketName,
				CheckName:  "Cross-Region Replication",
				Status:     models.StatusSkipped,
				Detail:     "Skipped: enable versioning first",
			})
		}
	}
	return results, nil
}

func checkS3PublicAccessBlock(ctx context.Context, s3Client *s3.Client, bucketName string) models.AuditResult {
	result := models.AuditResult{Service: "S3", ResourceID: bucketName, CheckName: "Public Access Block"}
	output, err := s3Client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
		Bucket: &bucketName,
	})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "NoSuchPublicAccessBlockConfiguration") {
			result.Status = models.StatusFail
			result.Detail = "Public Access Block is NOT configured — bucket may be publicly accessible"
		} else {
			result.Status = models.StatusFail
			result.Detail = fmt.Sprintf("API error: %v", err)
		}
		return result
	}
	cfg := output.PublicAccessBlockConfiguration
	if cfg != nil &&
		cfg.BlockPublicAcls != nil && *cfg.BlockPublicAcls &&
		cfg.BlockPublicPolicy != nil && *cfg.BlockPublicPolicy &&
		cfg.IgnorePublicAcls != nil && *cfg.IgnorePublicAcls &&
		cfg.RestrictPublicBuckets != nil && *cfg.RestrictPublicBuckets {
		result.Status = models.StatusPass
		result.Detail = "All 4 Public Access Block settings are ENABLED"
	} else {
		result.Status = models.StatusFail
		result.Detail = "One or more Public Access Block settings are DISABLED — risk of data exposure"
	}
	return result
}

func checkS3Versioning(ctx context.Context, s3Client *s3.Client, bucketName string) models.AuditResult {
	result := models.AuditResult{Service: "S3", ResourceID: bucketName, CheckName: "Versioning"}
	output, err := s3Client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{Bucket: &bucketName})
	if err != nil {
		result.Status = models.StatusFail
		result.Detail = fmt.Sprintf("API error: %v", err)
		return result
	}
	if output.Status == types.BucketVersioningStatusEnabled {
		result.Status = models.StatusPass
		result.Detail = "Versioning is ENABLED"
	} else if output.Status == types.BucketVersioningStatusSuspended {
		result.Status = models.StatusWarning
		result.Detail = "Versioning is SUSPENDED"
	} else {
		result.Status = models.StatusFail
		result.Detail = "Versioning is DISABLED"
	}
	return result
}

func checkS3CRR(ctx context.Context, s3Client *s3.Client, bucketName, fallbackRegion string) models.AuditResult {
	result := models.AuditResult{Service: "S3", ResourceID: bucketName, CheckName: "Cross-Region Replication"}
	output, err := s3Client.GetBucketReplication(ctx, &s3.GetBucketReplicationInput{Bucket: &bucketName})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "ReplicationConfigurationNotFoundError") || strings.Contains(errMsg, "NoSuchReplicationConfiguration") {
			result.Status = models.StatusFail
			result.Detail = fmt.Sprintf("No CRR configured to %s", fallbackRegion)
		} else {
			result.Status = models.StatusFail
			result.Detail = fmt.Sprintf("API error: %v", err)
		}
		return result
	}
	for _, rule := range output.ReplicationConfiguration.Rules {
		if rule.Status == types.ReplicationRuleStatusDisabled {
			continue
		}
		if rule.Destination != nil && rule.Destination.Bucket != nil {
			if strings.Contains(*rule.Destination.Bucket, fallbackRegion) {
				result.Status = models.StatusPass
				result.Detail = fmt.Sprintf("Active CRR rule to %s", fallbackRegion)
				return result
			}
		}
	}
	result.Status = models.StatusFail
	result.Detail = fmt.Sprintf("No active CRR rule targeting %s", fallbackRegion)
	return result
}
