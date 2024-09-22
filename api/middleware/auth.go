package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type contextKey string
var AuthenticationKey = contextKey("Authentication")

// AuthMiddleware creates middleware to handle ZITADEL token introspection
func AuthMiddleware(provider rs.ResourceServer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the token exists and is valid
			ok, token := checkToken(w, r)
			if !ok {
				return
			}

			// Introspect the token with the ZITADEL provider
			resp, err := rs.Introspect[*oidc.IntrospectionResponse](r.Context(), provider, token)
			if err != nil || !resp.Active {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), AuthenticationKey, resp)

			// If the token is valid, pass to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper function to extract the token from the Authorization header
func checkToken(w http.ResponseWriter, r *http.Request) (bool, string) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return false, ""
	}
	if !strings.HasPrefix(auth, oidc.PrefixBearer) {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return false, ""
	}
	return true, strings.TrimPrefix(auth, oidc.PrefixBearer)
}
