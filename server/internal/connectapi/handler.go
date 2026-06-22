package connectapi

import (
	"context"
	"net/http"
	"strings"

	"connectrpc.com/connect"

	"github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1/flashcardv1connect"
	"github.com/Shivam583-hue/trueflashcard/server/internal/auth"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
	"github.com/Shivam583-hue/trueflashcard/server/internal/service"
)

const sessionCookieName = "session"

func NewHandler(q dbgen.Querier, sessions *auth.SessionManager, appURL string) http.Handler {
	opts := connect.WithInterceptors(authInterceptor(sessions))

	mux := http.NewServeMux()
	mux.Handle(flashcardv1connect.NewFolderServiceHandler(folderAPI{service.NewFolderService(q)}, opts))
	mux.Handle(flashcardv1connect.NewDeckServiceHandler(deckAPI{service.NewDeckService(q)}, opts))
	mux.Handle(flashcardv1connect.NewFlashcardServiceHandler(flashcardAPI{service.NewFlashcardService(q)}, opts))

	return withCORS(mux, appURL)
}

func authInterceptor(sessions *auth.SessionManager) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if token := tokenFromHeader(req.Header()); token != "" {
				if userID, err := sessions.Verify(token); err == nil {
					ctx = auth.WithUserID(ctx, userID)
				}
			}
			return next(ctx, req)
		}
	}
}

func tokenFromHeader(h http.Header) string {
	if authz := h.Get("Authorization"); authz != "" {
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
			return strings.TrimSpace(parts[1])
		}
	}
	if cookie, err := (&http.Request{Header: h}).Cookie(sessionCookieName); err == nil {
		return cookie.Value
	}
	return ""
}

func withCORS(next http.Handler, appURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" && origin == appURL {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms, Authorization, X-Grpc-Web, X-User-Agent")
			w.Header().Set("Access-Control-Expose-Headers", "Grpc-Status, Grpc-Message")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
