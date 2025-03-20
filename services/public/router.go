package public

import (
	"context"
	"net/http"
	"strings"

	v1 "github.com/BATUCHKA/real-estate-back/services/public/v1"
	"github.com/BATUCHKA/real-estate-back/util"
	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router) {
	r.Use(baseRoute)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/project", v1.ProjectListGet)
		r.Get("/project/{id}", v1.ProjectGetByID)
	})
}

func baseRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		bearerToken := strings.TrimPrefix(authorization, "Bearer ")
		user, session, _, err := util.ParseUserSession(bearerToken)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), util.SessionKey, session))
		r = r.WithContext(context.WithValue(r.Context(), util.SessionUserKey, user))
		next.ServeHTTP(w, r)
	})
}
