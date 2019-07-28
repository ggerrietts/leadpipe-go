package collect

import (
	"bufio"
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-uuid"

	"github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/internal/kafka"
)

const (
	// CookieName is the name of the cookie to check for a visitor ID.
	CookieName = "visitor_id"
	// QueryName is the name of the query parameter used to convey the payload
	QueryName = "q"
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

// retrieveVisitorIDFromCookie does exactly that.
func retrieveVisitorIDFromCookie(r *http.Request) (string, error) {
	// retrieve visitor ID from cookies
	var visitor string
	cookie, err := r.Cookie(CookieName)
	if err == nil {
		return cookie.Value, nil
	}

	visitor, err = uuid.GenerateUUID()
	if err != nil {
		logrus.WithError(err).Error("unable to generate UUID")
		return "", err
	}

	r.AddCookie(&http.Cookie{
		Name:     CookieName,
		Value:    visitor,
		HttpOnly: true,
	})
	return visitor, nil
}

// retrieveDataFromURL
func retrieveDataFromQueryString(r *http.Request) string {
	payload := r.URL.Query().Get(QueryName)
	return payload
}

func buildImage() []byte {
	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	draw.Draw(m, m.Bounds(), image.Transparent, image.ZP, draw.Src)

	var b bytes.Buffer
	err := png.Encode(bufio.NewWriter(&b), m)
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
		visitor, err := retrieveVisitorIDFromCookie(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\": \"Internal service error\"}"))
			return
		}

		payload := retrieveDataFromQueryString(r)

		producer.Send(r.Context(), visitor, payload)

		sendImage(w)
	}
}
