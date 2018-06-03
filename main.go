package main

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/davyj0nes/s3-region-stats/sorter"
)

const padding = 3

// need to be here so they can be accessed by
var (
	bucketRegions []string
	wg            sync.WaitGroup
)

func main() {
	start := time.Now()

	regionStats := getRegionStats()
	textOutput(regionStats)

	fmt.Printf("\n%.2fs elapsed\n", time.Since(start).Seconds())
}

func getRegionStats() sorter.StatList {
	s3 := InitialiseClient()
	buckets := GetAllBuckets(s3)

	// Add number of buckets to the waitgroup
	wg.Add(len(buckets))
	for _, bucket := range buckets {
		go getBucketRegion(s3, bucket, &wg)
	}

	// block until all goroutines have completed
	wg.Wait()

	regionStats := make(map[string]int)
	for _, item := range bucketRegions {
		_, exist := regionStats[item]
		if exist {
			regionStats[item]++
		} else {
			regionStats[item] = 1
		}
	}

	return sorter.Sorter(regionStats)
}

// GetBucketRegion gets the region of a bucket and sends to to a channel
func getBucketRegion(svc *s3.S3, bucket string, wg *sync.WaitGroup) {
	// schedule a call to Done so we can decrement the wg list
	defer wg.Done()

	// input params for S3 call
	input := &s3.GetBucketLocationInput{
		Bucket: &bucket,
	}

	// create S3 request object
	req := svc.GetBucketLocationRequest(input)

	// make API call
	output, err := req.Send()
	if err != nil {
		panic("Error making API request, " + err.Error())
	}

	// update global variable
	bucketRegions = append(bucketRegions, string(output.LocationConstraint))
}

func textOutput(regionStats sorter.StatList) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	total := 0

	fmt.Fprintln(w, "Region\tCount\t")
	for _, stat := range regionStats {
		if stat.Key == "" {
			stat.Key = "No Region"
		}
		total += stat.Value
		fmt.Fprintf(w, "%s\t %d\t\n", stat.Key, stat.Value)
	}

	fmt.Fprintf(w, "TOTAL\t %d\t\n", total)

	w.Flush()
}

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
	// input params for S3 call
	input := &s3.ListBucketsInput{}

	// create S3 request object
	req := svc.ListBucketsRequest(input)

	// make API call
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
