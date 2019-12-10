package sources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/s3"
	lru "github.com/hashicorp/golang-lru"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
)

var (
	// Bucket name -> region
	bucketCache *lru.ARCCache

	// region -> S3 client
	clientCache map[string]*s3.S3
)

func init() {
	var err error
	clientCache = map[string]*s3.S3{}

	bucketCache, err = lru.NewARC(1000)
	if err != nil {
		panic("Failed to create bucket cache")
	}
}

// getS3Client Fetches S3 client with permissions to read data from the given account
func getS3Client(s3Bucket string) (*s3.S3, error) {
	var err error
	bucketRegion, ok := bucketCache.Get(s3Bucket)
	if !ok {
		zap.L().Info("bucket region was not cached, fetching it", zap.String("bucket", s3Bucket))
		bucketRegion, err = getBucketRegion(s3Bucket)
		if err != nil {
			return nil, err
		}
		bucketCache.Add(s3Bucket, bucketRegion)
	}

	zap.L().Debug("found bucket region", zap.Any("region", bucketRegion))

	client, ok := clientCache[bucketRegion.(string)]
	if !ok {
		zap.L().Info("s3 client was not cached, creating it")
		client = s3.New(common.Session, aws.NewConfig().
			WithRegion(bucketRegion.(string)))
		clientCache[bucketRegion.(string)] = client
	}
	return client, nil
}

func getBucketRegion(s3Bucket string) (string, error) {
	zap.L().Info("searching bucket region",
		zap.String("bucket", s3Bucket))

	locationDiscoveryClient := s3.New(common.Session)
	input := &s3.GetBucketLocationInput{Bucket: aws.String(s3Bucket)}
	location, err := locationDiscoveryClient.GetBucketLocation(input)
	if err != nil {
		zap.L().Warn("failed to find bucket region", zap.Error(err))
		return "", err
	}

	// Method may return nil if region is us-east-1,https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLocation.html
	// and https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region
	if location.LocationConstraint == nil {
		return endpoints.UsEast1RegionID, nil
	}
	return *location.LocationConstraint, nil
}
