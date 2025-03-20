package api

import (
	"context"
	"net/http"
	"strings"

	v1 "github.com/BATUCHKA/real-estate-back/services/api/v1"
	"github.com/BATUCHKA/real-estate-back/util"
	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router) {
	r.Use(baseRoute)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/ifno", v1.SettingsInfoGet)
		r.Put("/info", v1.SettingsInfoPut)
		r.Post("/project/{id}", v1.ProjectCreate)
		r.Put("/project/{id}", v1.ProjectUpdateByID)
		r.Delete("/project/{id}", v1.ProjectDeleteByID)
	})
}

func baseRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// whiteLabel := [...]string{"/v1/auth/login", "/v1/auth/signup", "/v1/auth/confirm/email/*", "/v1/auth/password/*", "/v1/auth/otp/confirm"}
		whiteLabel := [...]string{"/v1/auth/login", "/v1/auth/signup", "/v1/auth/password/*"}
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
		user, session, _, err := util.ParseUserSession(bearerToken)
		if err != nil {
			util.JsonErrorResponse("Not Authenticated.").WithErrorCode(401).Write(w)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), util.SessionKey, session))
		r = r.WithContext(context.WithValue(r.Context(), util.SessionUserKey, user))
		next.ServeHTTP(w, r)
	})
}
