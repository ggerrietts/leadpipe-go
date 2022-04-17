package collect

import (
	"encoding/json"
	"github.com/ggerrietts/leadpipe-go/pkg/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"

	"time"

	"github.com/ggerrietts/leadpipe-go/pkg/kafka"
)

var imageBytes []byte

// Server is the collection server
type Server struct {
	producer *kafka.Producer
	address  string
}

// NewServer generates a new Server
func NewServer(producer *kafka.Producer, address string) *Server {
	return &Server{
		producer: producer,
		address:  address,
	}
}

// Serve sets up the http handlers and serves forever
func (s *Server) Serve() {
	http.HandleFunc("/", s.ProcessPost)
	logrus.Fatal(http.ListenAndServe(s.address, nil))
}

func (s *Server) ProcessPost(w http.ResponseWriter, r *http.Request) {
	var val map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&val)
	if err != nil {
		logrus.WithError(err).Error("failure parsing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logrus.WithField("payload", val).Info("received")
	payload, err := structpb.NewStruct(val)
	if err != nil {
		logrus.WithError(err).Error("failure converting to protobuf")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	evt := pb.LeadpipeEvent{
		UserID:      0,
		ReceiveTime: timestamppb.New(time.Now()),
		Payload:     payload,
	}
	s.producer.Send(r.Context(), &evt)

	w.WriteHeader(http.StatusOK)
}
