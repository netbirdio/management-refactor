package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/netbirdio/netbird/formatter/hook"
	"github.com/netbirdio/netbird/management/server/auth"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
	"github.com/netbirdio/netbird/management/server/http/middleware/bypass"
)

// AuthMiddleware middleware to verify personal access tokens (PAT) and JWT tokens
type AuthMiddleware struct {
	authManager auth.Manager
}

// NewAuthMiddleware instance constructor
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

// Handler method of the middleware which authenticates a user either by JWT claims or by PAT
func (m *AuthMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//nolint
		ctx := context.WithValue(r.Context(), hook.ExecutionContextKey, hook.HTTPSource)

		reqID := uuid.New().String()
		//nolint
		ctx = context.WithValue(ctx, nbcontext.RequestIDKey, reqID)

		if bypass.ShouldBypass(r.URL.Path, h, w, r) {
			return
		}

		auth := strings.Split(r.Header.Get("Authorization"), " ")

		userID := "forbiddenUser"
		switch auth[1] {
		case "allowed":
			userID = "allowedUser"

		}

		userAuth := nbcontext.UserAuth{
			AccountId:      "accountID",
			Domain:         "",
			DomainCategory: "",
			Invited:        false,
			IsChild:        false,
			UserId:         userID,
			LastLogin:      time.Time{},
			Groups:         nil,
			IsPAT:          false,
		}

		request := nbcontext.SetUserAuthInRequest(r, userAuth)
		h.ServeHTTP(w, request)
	})
}
