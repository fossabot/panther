package awslogs

import (
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
)

// AuroraMySQLAudit is an RDS Aurora audit log which contains context around database calls.
// Reference: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Auditing.html
type AuroraMySQLAudit struct {
	Timestamp    *time.Time `json:"timestamp,omitempty"`
	ServerHost   *string    `json:"serverHost,omitempty"`
	Username     *string    `json:"username,omitempty"`
	Host         *string    `json:"host,omitempty"`
	ConnectionID *int       `json:"connectionId,omitempty"`
	QueryID      *int       `json:"queryId,omitempty"`
	Operation    *string    `json:"operation,omitempty" validate:"oneof=CONNECT QUERY READ WRITE CREATE ALTER RENAME DROP"`
	Database     *string    `json:"database,omitempty"`
	Object       *string    `json:"object,omitempty"`
	RetCode      *int       `json:"retCode,omitempty"`
}

// AuroraMySQLAuditParser parses AWS Aurora MySQL Audit logs
type AuroraMySQLAuditParser struct{}

// Parse returns the parsed events or nil if parsing failed
func (p *AuroraMySQLAuditParser) Parse(log string) []interface{} {
	reader := csv.NewReader(strings.NewReader(log))
	records, err := reader.ReadAll()
	if err != nil {
		zap.L().Debug("failed to parse the log as csv")
		return nil
	}

	record := records[0]

	timestampUnixMillis, err := strconv.ParseInt(record[0], 0, 64)
	if err != nil {
		return nil
	}

	// If there are ',' in the "object" field, CSV reader will split it to multiple fields
	// We are concatenating them to re-create the field
	objectString := strings.Join(record[8:len(record)-1], ",")

	event := &AuroraMySQLAudit{
		Timestamp:    aws.Time(time.Unix(timestampUnixMillis/1000000, timestampUnixMillis%1000000*1000).In(time.UTC)),
		ServerHost:   csvStringToPointer(record[1]),
		Username:     csvStringToPointer(record[2]),
		Host:         csvStringToPointer(record[3]),
		ConnectionID: csvStringToIntPointer(record[4]),
		QueryID:      csvStringToIntPointer(record[5]),
		Operation:    csvStringToPointer(record[6]),
		Database:     csvStringToPointer(record[7]),
		Object:       csvStringToPointer(objectString),
		RetCode:      csvStringToIntPointer(record[len(record)-1]),
	}
	if err := parsers.Validator.Struct(event); err != nil {
		zap.L().Debug("failed to validate log", zap.Error(err))
		return nil
	}

	return []interface{}{event}
}

// LogType returns the log type supported by this parser
func (p *AuroraMySQLAuditParser) LogType() string {
	return "AWS.AuroraMySQLAudit"
}
