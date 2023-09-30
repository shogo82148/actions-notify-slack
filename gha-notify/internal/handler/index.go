package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"html/template"
	"net/http"
	"strconv"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
)

const sessionCookieName = "gha-notify-session-id"

//go:embed index.html
var indexTemplate string

type IndexHandler struct {
	cfg  *IndexHandlerConfig
	tmpl *template.Template
}

type IndexHandlerConfig struct {
	repository.SlackClientIDGetter
	repository.SessionGetter
	repository.SessionPutter
}

func NewIndexHandler(cfg *IndexHandlerConfig) (*IndexHandler, error) {
	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return nil, err
	}
	return &IndexHandler{
		cfg:  cfg,
		tmpl: tmpl,
	}, nil
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, header, err := h.render(ctx, r)
	if err != nil {
		handleError(r.Context(), w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Set-Cookie", header)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

type indexData struct {
	ClientID  string
	SessionID string
	State     string
	TeamID    string
	TeamName  string
}

func (h *IndexHandler) render(ctx context.Context, r *http.Request) ([]byte, string, error) {
	// get the session
	s, err := getSession(r, h.cfg.SessionGetter)
	if err != nil {
		return nil, "", err
	}
	if s.SessionID == "" {
		s.SessionID = newSessionID()
	}
	if s.State == "" {
		s.State = newState()
	}
	header, err := putSession(r, h.cfg.SessionPutter, s)
	if err != nil {
		return nil, "", err
	}

	// get the client ID
	clientID, err := h.cfg.GetSlackClientID(ctx, &repository.GetSlackClientIDInput{})
	if err != nil {
		return nil, "", err
	}

	data := &indexData{
		ClientID:  clientID.SlackClientID,
		SessionID: s.SessionID,
		State:     s.State,
		TeamID:    s.TeamID,
		TeamName:  s.TeamName,
	}
	buf := new(bytes.Buffer)
	if err := h.tmpl.Execute(buf, data); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), header, nil
}

type session struct {
	SessionID string
	State     string
	TeamID    string
	TeamName  string
}

func getSession(req *http.Request, getter repository.SessionGetter) (*session, error) {
	ctx := req.Context()
	id := getSessionID(req)
	if id == "" {
		return &session{
			SessionID: newSessionID(),
		}, nil
	}
	s, err := getter.GetSession(ctx, &repository.GetSessionInput{
		SessionID: id,
	})
	if err != nil {
		return nil, err
	}
	return &session{
		SessionID: s.SessionID,
		State:     s.State,
		TeamID:    s.TeamID,
		TeamName:  s.TeamName,
	}, nil
}

func putSession(req *http.Request, putter repository.SessionPutter, s *session) (header string, err error) {
	ctx := req.Context()
	_, err = putter.PutSession(ctx, &repository.PutSessionInput{
		SessionID: s.SessionID,
		State:     s.State,
		TeamID:    s.TeamID,
		TeamName:  s.TeamName,
	})
	if err != nil {
		return "", err
	}
	c := &http.Cookie{
		Name:     sessionCookieName,
		Value:    s.SessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return c.String(), nil
}

func newSessionID() string {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf[:])
}

func newState() string {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf[:])
}

func getSessionID(req *http.Request) string {
	for _, c := range req.Cookies() {
		if c.Name == sessionCookieName {
			return c.Value
		}
	}
	return ""
}
