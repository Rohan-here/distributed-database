package web

import (
	"distributeddb/db"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
)

// Server contains HTTP method handlers to be used
type Server struct {
	db         *db.Database
	shardIdx   int
	shardCount int
	addrs      map[int]string
}

// Creates new server instance with HTTP handlers
func NewServer(db *db.Database, shardIdx int, shardCount int, addrs map[int]string) *Server {
	return &Server{
		db:         db,
		shardIdx:   shardIdx,
		shardCount: shardCount,
		addrs:      addrs,
	}
}

func (s *Server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shardCount))
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.addrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d to %d\n", s.shardIdx, shard)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error redirecting the request = %v", err)
		return
	}

	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

// Handles read requests
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	shard := s.getShard(key)
	value, err := s.db.GetKey(key)

	if shard != s.shardIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "Value = %q, Error = %v, Shard= %d, address : %q", value, err, shard, s.addrs[shard])
}

// Handles write requests
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	log.Print(key)

	shard := s.getShard(key)
	if shard != s.shardIdx {
		s.redirect(shard, w, r)
		return
	}

	err := s.db.SetKey(key, []byte(value))

	fmt.Fprintf(w, "Error = %v, Shard= %d", err, shard)
}
