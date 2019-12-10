package awslogs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

func TestAuroraMySQLAuditLog(t *testing.T) {
	//nolint:lll
	log := "1572546356975302,db-instance-name,someuser,10.0.143.147,1688069,168876953,QUERY,testdb," +
		"'select `user_id` as `userId`, `address`, `type`, `access`, `ordinal`, `token`, `verified`, `organization_id` as `organizationId`, `expires_at` as `expiresAt`, `created_at` as `createdAt`, `updated_at` as `updatedAt` " +
		"from `address_verification` where `ordinal` = \\'primary\\' and `access` = \\'public\\' and `type` = \\'phoneNumber\\' and `verified` = true and `user_id` = \\'12345678-8a3b-4d3f-96a7-19cc4c58c25d\\'',0"

	expectedTime := time.Unix(1572546356, 975302000).In(time.UTC)
	expectedEvent := &AuroraMySQLAudit{
		Timestamp:    aws.Time(expectedTime),
		ServerHost:   aws.String("db-instance-name"),
		Username:     aws.String("someuser"),
		Host:         aws.String("10.0.143.147"),
		ConnectionID: aws.Int(1688069),
		QueryID:      aws.Int(168876953),
		Operation:    aws.String("QUERY"),
		Database:     aws.String("testdb"),
		//nolint:lll
		Object:  aws.String("'select `user_id` as `userId`, `address`, `type`, `access`, `ordinal`, `token`, `verified`, `organization_id` as `organizationId`, `expires_at` as `expiresAt`, `created_at` as `createdAt`, `updated_at` as `updatedAt` from `address_verification` where `ordinal` = \\'primary\\' and `access` = \\'public\\' and `type` = \\'phoneNumber\\' and `verified` = true and `user_id` = \\'12345678-8a3b-4d3f-96a7-19cc4c58c25d\\''"),
		RetCode: aws.Int(0),
	}
	parser := &AuroraMySQLAuditParser{}
	require.Equal(t, []interface{}{expectedEvent}, parser.Parse(log))
}

func TestAuroraMysqlAuditLogType(t *testing.T) {
	parser := &AuroraMySQLAuditParser{}
	require.Equal(t, "AWS.AuroraMySQLAudit", parser.LogType())
}
