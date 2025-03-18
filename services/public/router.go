package public

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
		// r.Get("/{short_key}", v1.SurveyLinkGet)
		// r.Get("/surveys", v1.SurveyList)
		// r.Get("/banks", v1.BanksGet)
		// r.Get("/surveys/{survey_id}", v1.SurveyGetByID)
		// // r.Get("/surveys/section/{section_id}", v1.SurveySectionGet)

		// r.Get("/image/{name}", v1.ImageDownload)
	})
}

func baseRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		bearerToken := strings.TrimPrefix(authorization, "Bearer ")
		user, session, err := util.ParseSession(bearerToken)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), util.SessionKey, session))
		r = r.WithContext(context.WithValue(r.Context(), util.SessionUserKey, user))
		next.ServeHTTP(w, r)
	})
}
