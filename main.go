package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/davyj0nes/s3-region-stats/awsapi"
	"github.com/davyj0nes/s3-region-stats/sorter"
)

const padding = 3

func main() {
	start := time.Now()

	regionStats := getRegionStats()
	textOutput(regionStats)

	fmt.Printf("\n%.2fs elapsed\n", time.Since(start).Seconds())
}

func getRegionStats() sorter.StatList {
	s3 := awsapi.InitialiseClient()
	bucketRegions := []string{}
	regionChannel := make(chan string)

	buckets := awsapi.GetAllBuckets(s3)
	// generate output to channel
	for _, bucket := range buckets {
		go awsapi.GetBucketRegion(s3, bucket, regionChannel)
	}

	// read from channel
	for range buckets {
		bucketRegions = append(bucketRegions, <-regionChannel)
	}

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
