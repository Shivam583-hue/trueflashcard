package httpauth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Shivam583-hue/trueflashcard/server/internal/auth"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
)

const (
	stateCookieName   = "oauth_state"
	sessionCookieName = "session"
	stateTTLSeconds   = 600
)

type Handler struct {
	oauth    *auth.GoogleOAuth
	sessions *auth.SessionManager
	users    dbgen.Querier
	appURL   string
	secure   bool
}

func NewHandler(oauth *auth.GoogleOAuth, sessions *auth.SessionManager, users dbgen.Querier) *Handler {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	return &Handler{
		oauth:    oauth,
		sessions: sessions,
		users:    users,
		appURL:   appURL,
		secure:   os.Getenv("COOKIE_SECURE") == "true",
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /auth/google/login", h.login)
	mux.HandleFunc("GET /auth/google/callback", h.callback)
	mux.HandleFunc("POST /auth/logout", h.logout)
	mux.HandleFunc("GET /auth/me", h.me)
	return h.withCORS(mux)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	state, err := randomState()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Path:     "/",
		MaxAge:   stateTTLSeconds,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, h.oauth.AuthCodeURL(state), http.StatusFound)
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie(stateCookieName)
	if err != nil || state.Value == "" || state.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	h.clearCookie(w, stateCookieName)

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	identity, err := h.oauth.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	user, err := h.users.UpsertUserByGoogleSubject(r.Context(), dbgen.UpsertUserByGoogleSubjectParams{
		GoogleSubject: identity.Subject,
		Email:         identity.Email,
		DisplayName:   identity.Name,
	})
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, expiresAt, err := h.sessions.Issue(user.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, h.appURL+"/home", http.StatusFound)
}

func (h *Handler) logout(w http.ResponseWriter, _ *http.Request) {
	h.clearCookie(w, sessionCookieName)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	userID, err := h.sessions.Verify(cookie.Value)
	if err != nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	user, err := h.users.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	writeJSON(w, map[string]string{
		"id":    user.ID.String(),
		"email": user.Email,
		"name":  user.DisplayName,
	})
}

func (h *Handler) clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin == h.appURL {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
