package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
)

const sessionCookieName = "gha-notify-session-id"

type IndexHandler struct {
	cfg *IndexHandlerConfig
}

type IndexHandlerConfig struct {
	repository.SlackClientIDGetter
	repository.SessionGetter
	repository.SessionPutter
}

func NewIndexHandler(cfg *IndexHandlerConfig) (*IndexHandler, error) {
	return &IndexHandler{
		cfg: cfg,
	}, nil
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, err := getSession(r, h.cfg.SessionGetter)
	if err != nil {
		handleError(r.Context(), w, err)
		return
	}
	if s.SessionID == "" {
		s.SessionID = newSessionID()
	}
	if s.State == "" {
		s.State = newState()
	}
	header, err := putSession(r, h.cfg.SessionPutter, s)
	if err != nil {
		handleError(r.Context(), w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Set-Cookie", header)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s.SessionID))
}

func (h *IndexHandler) render() ([]byte, error) {
	buf := new(bytes.Buffer)
	return buf.Bytes(), nil
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

// <a href="https://slack.com/oauth/v2/authorize?client_id=118051372210.5947954494951&scope=chat:write,commands&user_scope="><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcSet="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>
