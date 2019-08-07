package collect

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/internal/kafka"
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
	http.HandleFunc("/img.png", generateImageHandler(s.producer))
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)
	logrus.Fatal(http.ListenAndServe(s.address, nil))
}

func buildImage() []byte {
	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	draw.Draw(m, m.Bounds(), image.Transparent, image.ZP, draw.Src)

	var b bytes.Buffer
	err := png.Encode(&b, m)
	if err != nil {
		logrus.WithError(err).Fatal("unable to encode PNG")
	}

	return b.Bytes()
}

func sendImage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(imageBytes)))
	if _, err := w.Write(imageBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("failed to send image")
	}
}

func generateImageHandler(producer *kafka.Producer) http.HandlerFunc {
	imageBytes = buildImage()
	return func(w http.ResponseWriter, r *http.Request) {
		hit, err := buildPayload(r)
		if err != nil {
			logrus.WithError(err).Error("unable to build payload")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\": \"Internal service error\"}"))
			return
		}

		producer.Send(r.Context(), hit)

		sendImage(w)
	}
}
