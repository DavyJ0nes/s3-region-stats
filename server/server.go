package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/davyj0nes/s3-region-stats/awsapi"
)

// Server Object
type Server struct {
	router *http.ServeMux
}

// NewServer Initalises the Web Server
func NewServer() *Server {
	r := http.NewServeMux()

	s := &Server{
		router: r,
	}

	// Add Routes to server
	s.routes()

	return s
}

// ServeHTTP is needed to ensure handlers work as expected
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router.HandleFunc("/stats", s.StatsHandler())
}

// StatsHandler gets S3 Bucket stats and returns as JSON
func (s *Server) StatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println("Calling Stats")
		stats := awsapi.GetRegionStats()
		data, err := json.Marshal(stats)
		if err != nil {
			log.Printf("Could not encode info data:\n%v", err)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
