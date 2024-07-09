package main

import (
	"distributeddb/config"
	"distributeddb/db"
	"distributeddb/web"
	"flag"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
)

var (
	dbLocation = flag.String("db-location", "", "Path to the bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host listen address")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard      = flag.String("shard", "", "The name of the shard for the data")
)

func parseFlags() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatalf("Must Provide DB location")
	}

	if *shard == "" {
		log.Fatalf("Must Provide shard name")
	}
}

func main() {
	parseFlags()

	var c config.Config

	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("toml.DecodeFile failed: %v", err)
	}

	var shardCount int
	var shardIdx int = -1
	var addrs = make(map[int]string)

	for _, s := range c.Shards {
		addrs[s.Idx] = s.Address

		if s.Idx+1 > shardCount {
			shardCount = s.Idx + 1
		}

		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatal("Shard not found")
	}

	log.Printf("Shard: %s, ShardIdx: %d, ShardCount: %d", *shard, shardIdx, shardCount)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}

	defer close()

	srv := web.NewServer(db, shardIdx, shardCount, addrs)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
