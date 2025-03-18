package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	// v1 "github.com/BATUCHKA/real-estate-back/services/admin/v1"
	"github.com/BATUCHKA/real-estate-back/util"
)

func Route(r chi.Router) {
	r.Use(baseRoute)
	r.Route("/v1", func(r chi.Router) {
	})
}

func baseRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		whiteLabel := [...]string{"/v1/auth/login", "/v1/auth/send-otp", "/v1/auth/verify-otp", "/v1/auth/password-create", "/v1/auth/password-forgot", "/v1/auth/password-reset"}
		for _, value := range whiteLabel {
			if strings.HasSuffix(value, "*") && strings.HasPrefix(chi.RouteContext(r.Context()).RoutePath, strings.TrimSuffix(value, "*")) {
				next.ServeHTTP(w, r)
				return
			}
			if chi.RouteContext(r.Context()).RoutePath == value {
				next.ServeHTTP(w, r)
				return
			}
		}

		authorization := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			util.JsonErrorResponse("Bearer token not found.").WithErrorCode(400).Write(w)
			return
		}
		bearerToken := strings.TrimPrefix(authorization, "Bearer ")
		user, session, err := util.ParseSession(bearerToken)
		if err != nil {
			util.JsonErrorResponse("Not Authenticated.").WithErrorCode(401).Write401(w)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), util.SessionKey, session))
		r = r.WithContext(context.WithValue(r.Context(), util.SessionUserKey, user))
		next.ServeHTTP(w, r)
	})
}
