package awsglue

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

const (
	s3Prefix = "foo/"
)

func TestGlueMetadata_PartitionPrefix(t *testing.T) {
	var gm *GlueMetadata
	var expected string

	refTime := time.Date(2020, 1, 3, 1, 1, 1, 0, time.UTC)

	gm = &GlueMetadata{
		s3Prefix:     s3Prefix,
		timebin:      GlueTableHourly,
		timeUnpadded: false,
	}
	expected = "foo/year=2020/month=01/day=03/hour=01/"
	assert.Equal(t, expected, gm.PartitionPrefix(refTime))
	gm.timeUnpadded = true
	expected = "foo/year=2020/month=1/day=3/hour=1/"
	assert.Equal(t, expected, gm.PartitionPrefix(refTime))

	gm = &GlueMetadata{
		s3Prefix:     s3Prefix,
		timebin:      GlueTableDaily,
		timeUnpadded: false,
	}
	expected = "foo/year=2020/month=01/day=03/"
	assert.Equal(t, expected, gm.PartitionPrefix(refTime))
	gm.timeUnpadded = true
	expected = "foo/year=2020/month=1/day=3/"
	assert.Equal(t, expected, gm.PartitionPrefix(refTime))

	gm = &GlueMetadata{
		s3Prefix:     s3Prefix,
		timebin:      GlueTableMonthly,
		timeUnpadded: false,
	}
	expected = "foo/year=2020/month=01/"
	assert.Equal(t, expected, gm.PartitionPrefix(refTime))
	gm.timeUnpadded = true
	expected = "foo/year=2020/month=1/"
	assert.Equal(t, expected, gm.PartitionPrefix(refTime))
}

func TestGlueMetadata_PartitionValues(t *testing.T) {
	var gm *GlueMetadata
	var expected []*string

	refTime := time.Date(2020, 1, 3, 1, 1, 1, 0, time.UTC)

	gm = &GlueMetadata{
		s3Prefix:     s3Prefix,
		timebin:      GlueTableHourly,
		timeUnpadded: false,
	}
	expected = []*string{
		aws.String(fmt.Sprintf("%d", refTime.Year())),
		aws.String(fmt.Sprintf("%02d", refTime.Month())),
		aws.String(fmt.Sprintf("%02d", refTime.Day())),
		aws.String(fmt.Sprintf("%02d", refTime.Hour())),
	}
	assert.Equal(t, expected, gm.PartitionValues(refTime))
	gm.timeUnpadded = true
	expected = []*string{
		aws.String(fmt.Sprintf("%d", refTime.Year())),
		aws.String(fmt.Sprintf("%d", refTime.Month())),
		aws.String(fmt.Sprintf("%d", refTime.Day())),
		aws.String(fmt.Sprintf("%d", refTime.Hour())),
	}
	assert.Equal(t, expected, gm.PartitionValues(refTime))

	gm = &GlueMetadata{
		s3Prefix:     s3Prefix,
		timebin:      GlueTableDaily,
		timeUnpadded: false,
	}
	expected = []*string{
		aws.String(fmt.Sprintf("%d", refTime.Year())),
		aws.String(fmt.Sprintf("%02d", refTime.Month())),
		aws.String(fmt.Sprintf("%02d", refTime.Day())),
	}
	assert.Equal(t, expected, gm.PartitionValues(refTime))
	gm.timeUnpadded = true
	expected = []*string{
		aws.String(fmt.Sprintf("%d", refTime.Year())),
		aws.String(fmt.Sprintf("%d", refTime.Month())),
		aws.String(fmt.Sprintf("%d", refTime.Day())),
	}
	assert.Equal(t, expected, gm.PartitionValues(refTime))

	gm = &GlueMetadata{
		s3Prefix:     s3Prefix,
		timebin:      GlueTableMonthly,
		timeUnpadded: false,
	}
	expected = []*string{
		aws.String(fmt.Sprintf("%d", refTime.Year())),
		aws.String(fmt.Sprintf("%02d", refTime.Month())),
	}
	assert.Equal(t, expected, gm.PartitionValues(refTime))
	gm.timeUnpadded = true
	expected = []*string{
		aws.String(fmt.Sprintf("%d", refTime.Year())),
		aws.String(fmt.Sprintf("%d", refTime.Month())),
	}
	assert.Equal(t, expected, gm.PartitionValues(refTime))
}
