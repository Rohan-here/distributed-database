package main

import (
	"distributeddb/db"
	"distributeddb/web"
	"flag"
	"log"
	"net/http"
)

var (
	dbLocation = flag.String("db-location", "", "Path to the bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host listen address")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
)

func parseFlags() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatalf("Must Provide DB location")
	}
}

func main() {
	parseFlags()

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}

	defer close()

	srv := web.NewServer(db)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
