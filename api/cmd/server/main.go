package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
)

type Zoo struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Pref           string  `json:"pref"`
	City           string  `json:"city"`
	Status         string  `json:"status"`
	OfficialURL    string  `json:"officialUrl"`
	LastVerifiedAt string  `json:"lastVerifiedAt"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
	Note           string  `json:"note,omitempty"`
}

func main() {
	ctx := context.Background()

	projectID := mustEnv("FIREBASE_PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// webのサーバー（python http.server）が 5500 の想定
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// ブラウザが事前に送る OPTIONS（プリフライト）への応答
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Get("/zoos", func(w http.ResponseWriter, r *http.Request) {
		zoos, err := listZoos(r.Context(), client)
		if err != nil {
			httpError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, zoos)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}

func listZoos(ctx context.Context, client *firestore.Client) ([]Zoo, error) {
	// Firestore collection: zoos
	docs, err := client.Collection("zoos").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	out := make([]Zoo, 0, len(docs))
	for _, d := range docs {
		m := d.Data()

		z := Zoo{
			ID:             d.Ref.ID,
			Name:           asString(m["name"]),
			Pref:           asString(m["pref"]),
			City:           asString(m["city"]),
			Status:         asString(m["status"]),
			OfficialURL:    asString(m["officialUrl"]),
			LastVerifiedAt: asString(m["lastVerifiedAt"]),
			Note:           asString(m["note"]),
		}

		// geo: {lat, lng} を想定
		if geo, ok := m["geo"].(map[string]any); ok {
			z.Lat = asFloat(geo["lat"])
			z.Lng = asFloat(geo["lng"])
		} else {
			// もし lat/lng 直置きにした場合にも動くように保険
			z.Lat = asFloat(m["lat"])
			z.Lng = asFloat(m["lng"])
		}

		out = append(out, z)
	}
	return out, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(v)
}

func httpError(w http.ResponseWriter, err error, code int) {
	log.Printf("error: %v", err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": err.Error(),
	})
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env: %s", key)
	}
	return v
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func asFloat(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return 0
	}
}
