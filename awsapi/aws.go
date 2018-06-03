package awsapi

import (
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// InitialiseClient provides an AWS Service object to interact with the API
func InitialiseClient() *s3.S3 {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	svc := s3.New(cfg)

	return svc
}

// GetAllBuckets gets the names of all S3 buckets in an account
func GetAllBuckets(svc *s3.S3) []string {
	input := &s3.ListBucketsInput{}

	req := svc.ListBucketsRequest(input)

	output, err := req.Send()
	if err != nil {
		panic("Error making API request, " + err.Error())
	}

	var bucketNames []string
	for _, bucket := range output.Buckets {
		bucketNames = append(bucketNames, *bucket.Name)
	}

	return bucketNames
}

// GetBucketRegion gets the region of a bucket and sends to to a channel
func GetBucketRegion(svc *s3.S3, bucket string, c chan string) {
	input := &s3.GetBucketLocationInput{
		Bucket: &bucket,
	}

	req := svc.GetBucketLocationRequest(input)

	output, err := req.Send()
	if err != nil {
		panic("Error making API request, " + err.Error())
	}
	c <- string(output.LocationConstraint)
}
