package auditors

import (
	"context"
	"fmt"

	awsclients "github.com/Mohamed-M-Meth/aws-dr-audit/internal/aws"
	"github.com/Mohamed-M-Meth/aws-dr-audit/pkg/models"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func AuditRDS(ctx context.Context, clients *awsclients.Clients, fallbackRegion string) ([]models.AuditResult, error) {
	var results []models.AuditResult
	var marker *string
	var allInstances []rdsInstance
	for {
		output, err := clients.RDSPrimary.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{Marker: marker})
		if err != nil {
			return nil, fmt.Errorf("DescribeDBInstances failed: %w", err)
		}
		for _, db := range output.DBInstances {
			multiAZ := false
			if db.MultiAZ != nil {
				multiAZ = *db.MultiAZ
			}
			allInstances = append(allInstances, rdsInstance{
				id:      *db.DBInstanceIdentifier,
				engine:  *db.Engine,
				multiAZ: multiAZ,
				status:  *db.DBInstanceStatus,
			})
		}
		if output.Marker == nil {
			break
		}
		marker = output.Marker
	}
	if len(allInstances) == 0 {
		fmt.Println("  No RDS instances found.")
		return results, nil
	}
	fmt.Printf("  Found %d RDS instance(s). Checking each...\n\n", len(allInstances))
	for _, db := range allInstances {
		if db.status != "available" {
			results = append(results, models.AuditResult{
				Service:    "RDS",
				ResourceID: db.id,
				CheckName:  "Multi-AZ",
				Status:     models.StatusSkipped,
				Detail:     fmt.Sprintf("Skipped: status is '%s'", db.status),
			})
			continue
		}
		results = append(results, checkRDSMultiAZ(db))
		crossRegionResult, err := checkRDSCrossRegionSnapshots(ctx, clients.RDSFallback, db.id, fallbackRegion)
		if err != nil {
			results = append(results, models.AuditResult{
				Service:    "RDS",
				ResourceID: db.id,
				CheckName:  "Cross-Region Snapshots",
				Status:     models.StatusFail,
				Detail:     fmt.Sprintf("Error: %v", err),
			})
		} else {
			results = append(results, crossRegionResult)
		}
	}
	return results, nil
}

type rdsInstance struct {
	id      string
	engine  string
	multiAZ bool
	status  string
}

func checkRDSMultiAZ(db rdsInstance) models.AuditResult {
	result := models.AuditResult{Service: "RDS", ResourceID: db.id, CheckName: "Multi-AZ"}
	if db.multiAZ {
		result.Status = models.StatusPass
		result.Detail = "Multi-AZ is ENABLED"
	} else {
		result.Status = models.StatusFail
		result.Detail = "Multi-AZ is DISABLED"
	}
	return result
}

func checkRDSCrossRegionSnapshots(ctx context.Context, rdsClientFallback *rds.Client, instanceID, fallbackRegion string) (models.AuditResult, error) {
	result := models.AuditResult{Service: "RDS", ResourceID: instanceID, CheckName: "Cross-Region Snapshots"}
	output, err := rdsClientFallback.DescribeDBSnapshots(ctx, &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: &instanceID,
	})
	if err != nil {
		return result, fmt.Errorf("DescribeDBSnapshots in %s failed: %w", fallbackRegion, err)
	}
	for _, snap := range output.DBSnapshots {
		if snap.Status != nil && *snap.Status == "available" {
			result.Status = models.StatusPass
			result.Detail = fmt.Sprintf("Snapshot found in %s", fallbackRegion)
			return result, nil
		}
	}
	result.Status = models.StatusFail
	result.Detail = fmt.Sprintf("No snapshots found in %s", fallbackRegion)
	return result, nil
}
