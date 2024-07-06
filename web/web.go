package web

import (
	"distributeddb/db"
	"fmt"
	"net/http"
)

// Server contains HTTP method handlers to be used
type Server struct {
	db *db.Database
}

// Creates new server instance with HTTP handlers
func NewServer(db *db.Database) *Server {
	return &Server{db: db}
}

// Handles read requests
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	value, err := s.db.GetKey(key)

	fmt.Fprintf(w, "Value = %q, error = %v", value, err)
}

// Handles write requests
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err := s.db.SetKey(key, []byte(value))

	fmt.Fprintf(w, "Error = %v", err)
}
