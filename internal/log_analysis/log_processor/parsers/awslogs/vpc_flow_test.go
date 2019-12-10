package awslogs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestVpcFlowLog(t *testing.T) {
	parser := &VPCFlowParser{}

	log := "2 348372346321 eni-00184058652e5a320 52.119.169.95 172.31.20.31 443 48316 6 19 7119 1573642242 1573642284 ACCEPT OK"

	expectedStartTime := time.Unix(1573642242, 0).In(time.UTC)
	expectedEndTime := time.Unix(1573642284, 0).In(time.UTC)
	expectedEvent := &VPCFlow{
		Action:      aws.String("ACCEPT"),
		Account:     aws.String("348372346321"),
		Bytes:       aws.Int(7119),
		Dstaddr:     aws.String("172.31.20.31"),
		DstPort:     aws.Int(48316),
		End:         aws.Time(expectedEndTime),
		InterfaceID: aws.String("eni-00184058652e5a320"),
		LogStatus:   aws.String("OK"),
		Packets:     aws.Int(19),
		Protocol:    aws.Int(6),
		SourceAddr:  aws.String("52.119.169.95"),
		SrcPort:     aws.Int(443),
		Start:       aws.Time(expectedStartTime),
		Version:     aws.Int(2),
	}

	result := parser.Parse(log)
	// Testify can fail when comparing time.Time, since it compares the strings converted to local time
	if !cmp.Equal([]interface{}{expectedEvent}, result) {
		t.Fail()
	}
}

func TestVpcFlowLogNoData(t *testing.T) {
	parser := &VPCFlowParser{}

	log := "2 unknown eni-0608192d5c498fbcd - - - - - - - 1538696170 1538696308 - NODATA"

	expectedStartTime := time.Unix(1538696170, 0).In(time.UTC)
	expectedEndTime := time.Unix(1538696308, 0).In(time.UTC)
	expectedEvent := &VPCFlow{
		Version:     aws.Int(2),
		InterfaceID: aws.String("eni-0608192d5c498fbcd"),
		Start:       aws.Time(expectedStartTime),
		End:         aws.Time(expectedEndTime),
		LogStatus:   aws.String("NODATA"),
	}

	result := parser.Parse(log)
	// Testify can fail when comparing time.Time, since it compares the strings converted to local time
	if !cmp.Equal([]interface{}{expectedEvent}, result) {
		t.Fail()
	}
}

func TestVpcFlowLogType(t *testing.T) {
	parser := &VPCFlowParser{}
	require.Equal(t, "AWS.VPCFlow", parser.LogType())
}
