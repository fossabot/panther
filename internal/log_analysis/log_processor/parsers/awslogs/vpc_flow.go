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

// VPCFlow is a VPC NetFlow log, which is a layer 3 representation of network traffic in EC2.
// Log format & samples can be seen here: https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs-records-examples.html
type VPCFlow struct {
	Version     *int       `json:"version,omitempty" validate:"required"`
	Account     *string    `json:"account,omitempty" validate:"omitempty,len=12,numeric"`
	InterfaceID *string    `json:"interfaceId,omitempty"`
	SourceAddr  *string    `json:"sourceAddr,omitempty"`
	Dstaddr     *string    `json:"dstAddr,omitempty"`
	SrcPort     *int       `json:"srcPort,omitempty" validate:"omitempty,min=0,max=65535"`
	DstPort     *int       `json:"destPort,omitempty" validate:"omitempty,min=0,max=65535"`
	Protocol    *int       `json:"protocol,omitempty"`
	Packets     *int       `json:"packets,omitempty"`
	Bytes       *int       `json:"bytes,omitempty"`
	Start       *time.Time `json:"start,omitempty" validate:"required"`
	End         *time.Time `json:"end,omitempty" validate:"required"`
	Action      *string    `json:"action,omitempty" validate:"omitempty,oneof=ACCEPT REJECT"`
	LogStatus   *string    `json:"status,omitempty" validate:"oneof=OK NODATA SKIPDATA"`
}

// VPCFlowParser parses AWS VPC Flow Parser logs
type VPCFlowParser struct{}

// Parse returns the parsed events or nil if parsing failed
func (p *VPCFlowParser) Parse(log string) []interface{} {
	reader := csv.NewReader(strings.NewReader(log))
	reader.Comma = ' '

	records, err := reader.ReadAll()
	if err != nil {
		zap.L().Debug("failed to parse the log as csv")
		return nil
	}

	record := records[0]

	var account *string = nil
	if record[1] != "-" && record[1] != "unknown" {
		account = &record[1]
	}

	startTimeUnix, err := strconv.Atoi(record[10])
	if err != nil {
		return nil
	}
	endTimeUnix, err := strconv.Atoi(record[11])
	if err != nil {
		return nil
	}

	startTime := time.Unix(int64(startTimeUnix), 0)
	endTime := time.Unix(int64(endTimeUnix), 0)

	event := &VPCFlow{
		Version:     csvStringToIntPointer(record[0]),
		Account:     account,
		InterfaceID: csvStringToPointer(record[2]),
		SourceAddr:  csvStringToPointer(record[3]),
		Dstaddr:     csvStringToPointer(record[4]),
		SrcPort:     csvStringToIntPointer(record[5]),
		DstPort:     csvStringToIntPointer(record[6]),
		Protocol:    csvStringToIntPointer(record[7]),
		Packets:     csvStringToIntPointer(record[8]),
		Bytes:       csvStringToIntPointer(record[9]),
		Start:       aws.Time(startTime.In(time.UTC)),
		End:         aws.Time(endTime.In(time.UTC)),
		Action:      csvStringToPointer(record[12]),
		LogStatus:   csvStringToPointer(record[13]),
	}

	if err := parsers.Validator.Struct(event); err != nil {
		zap.L().Debug("failed to validate log", zap.Error(err))
		return nil
	}

	return []interface{}{event}
}

// LogType returns the log type supported by this parser
func (p *VPCFlowParser) LogType() string {
	return "AWS.VPCFlow"
}
