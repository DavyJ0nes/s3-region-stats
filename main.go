package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"text/tabwriter"
	"time"

	"github.com/davyj0nes/s3-region-stats/awsapi"
	"github.com/davyj0nes/s3-region-stats/server"
	"github.com/davyj0nes/s3-region-stats/sorter"
)

const padding = 3

func main() {
	enableTrace := flag.Bool("trace", false, "Enable Tracing or not")
	runServer := flag.Bool("server", false, "Run App as Web Service")
	flag.Parse()

	if *enableTrace {
		trace.Start(os.Stderr)
		defer trace.Stop()
	}

	start := time.Now()

	if *runServer {
		srv := server.NewServer()
		log.Println("Starting Server")
		log.Fatal(http.ListenAndServe("0.0.0.0:8008", srv))
	} else {
		regionStats := awsapi.GetRegionStats()
		textOutput(sorter.Sorter(regionStats))
		fmt.Printf("\n%.2fs elapsed\n", time.Since(start).Seconds())
	}
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
