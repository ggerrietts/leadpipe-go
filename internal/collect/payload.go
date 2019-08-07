package collect

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/ggerrietts/leadpipe-go/internal/pb"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
)

const (
	// CookieName is the name of the cookie to check for a visitor ID.
	CookieName = "visitor_id"
	// QueryName is the name of the query parameter used to convey the payload
	QueryName = "q"
)

// HitJSON is a decoder struct for the incoming JSON on a hit
type HitJSON struct {
	Token string `json:"token"`
	User  string `json:"user"`
	Event string `json:"event"`
}

// getVisitorID retrieves the visitor ID from the visitor ID cookie. If no
// visitor ID is found, it creates a new one and sets that into the request.
func getVisitorID(r *http.Request) (string, error) {
	// retrieve visitor ID from cookies
	var visitor string
	cookie, err := r.Cookie(CookieName)
	if err == nil {
		return cookie.Value, nil
	}

	visitor, err = uuid.GenerateUUID()
	if err != nil {
		return "", errors.Wrap(err, "getVisitorID")
	}

	r.AddCookie(&http.Cookie{
		Name:     CookieName,
		Value:    visitor,
		HttpOnly: true,
	})
	return visitor, nil
}

// decodeQuery retrieves the query from the request, un-base64's it,
// then parses the JSON
func decodeQuery(r *http.Request) (*HitJSON, error) {
	var data HitJSON
	payload := r.URL.Query().Get(QueryName)
	jsonBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, errors.Wrap(err, "decodeQuery - base64")
	}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, errors.Wrap(err, "decodeQuery - json")
	}
	return &data, nil
}

// HitHeaders gathers the information retrieved from headers
type HitHeaders struct {
	Page  string
	Addr  string
	Agent string
}

// getHeaders retrieves hit data from the request headers
func getHeaders(r *http.Request) *HitHeaders {
	return &HitHeaders{
		Page:  r.Referer(),
		Addr:  r.RemoteAddr,
		Agent: r.Header.Get("User-Agent"),
	}
}

// buildPayload is the main entrypoint to this file. It coordinates
// the building of a payload.
func buildPayload(r *http.Request) (*pb.Hit, error) {
	visitorID, err := getVisitorID(r)
	if err != nil {
		return nil, errors.Wrap(err, "buildPayload")
	}
	queryData, err := decodeQuery(r)
	if err != nil {
		return nil, errors.Wrap(err, "buildPayload")
	}
	headerData := getHeaders(r)

	return &pb.Hit{
		VisitorID: visitorID,
		Token:     queryData.Token,
		User:      queryData.User,
		Event:     queryData.Event,
		Page:      headerData.Page,
		Addr:      headerData.Addr,
		Agent:     headerData.Agent,
	}, nil
}
