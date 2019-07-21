package collect

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/internal/kafka"
)

// Server is the collection server
type Server struct {
	producer kafka.Producer
}

// NewServer generates a new Server
func NewServer(producer kafka.Producer) *Server {
	return &Server{
		producer: producer,
	}
}

// Serve sets up the http handlers and serves forever
func (s *Server) Serve() {
	http.Handle("/", generateImageHandler(s.producer))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateImageHandler(producer kafka.Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("preparing to produce")
		producer.Send("1234", "hello sweet world")
		log.Debug("preparing to write")
		w.Write([]byte("righteous"))
	}
}
