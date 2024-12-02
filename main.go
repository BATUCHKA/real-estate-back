package main

import (
	// "context"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var translations map[string]map[string]string
var translationsMutex sync.RWMutex

func main() {
	// database connection initialize
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	port := os.Getenv("SERVER_PORT")
	log.Println("chi service started on port", port)

	fs := http.FileServer(http.Dir("data"))
	r.Handle("/img/*", http.StripPrefix("/img/", fs))

	r.Mount("/", router())

	err = http.ListenAndServe(fmt.Sprint(":", port), r)
	if err != nil {
		log.Fatal(err)
	}
}

// func loadTranslations(filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Printf("Error opening file: %s\n", err)
// 		return
// 	}
// 	defer file.Close()

// 	decoder := json.NewDecoder(file)
// 	var messages map[string]string
// 	if err := decoder.Decode(&messages); err != nil {
// 		fmt.Printf("Error decoding JSON: %s\n", err)
// 		return
// 	}

// 	lang := filename[:2]
// 	translationsMutex.Lock()
// 	translations[lang] = messages
// 	translationsMutex.Unlock()
// }

func router() http.Handler {
	r := chi.NewRouter()
	r.Use(baseRoute)
	r.Route("/api", api.Route)
	r.Route("/admin/api", admin.Route)
	r.Route("/p/api", public.Route)
	return r
}

func baseRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
