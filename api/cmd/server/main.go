package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
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

type Koala struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Sex          string `json:"sex"`
	BirthDate    string `json:"birthDate"`
	Status       string `json:"status"`
	CurrentZooID string `json:"currentZooId"`
	MotherID     string `json:"motherId"`
	FatherID     string `json:"fatherId"`
}

type PedigreeNode struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	BirthDate string `json:"birthDate"`
	Note      string `json:"note,omitempty"`
}

type PedigreeEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"` // "mother" | "father"
}

type PedigreeResponse struct {
	RootID string         `json:"rootId"`
	Nodes  []PedigreeNode `json:"nodes"`
	Edges  []PedigreeEdge `json:"edges"`
}

type serverConfig struct {
	ProjectID   string
	Port        string
	ServeStatic bool
	WebDir      string
}

func main() {
	ctx := context.Background()

	cfg := loadServerConfig()
	client, err := firestore.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	r := chi.NewRouter()

	api := chi.NewRouter()
	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	api.Get("/zoos", func(w http.ResponseWriter, r *http.Request) {
		zoos, err := listZoos(r.Context(), client)
		if err != nil {
			httpError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, zoos)
	})
	api.Get("/zoos/{zooId}/koalas", func(w http.ResponseWriter, r *http.Request) {
		zooId := chi.URLParam(r, "zooId")
		status := r.URL.Query().Get("status")
		if status == "" {
			status = "alive"
		}

		koalas, err := listKoalasByZoo(r.Context(), client, zooId, status)
		if err != nil {
			httpError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, koalas)
	})

	api.Get("/koalas/{koalaId}/pedigree", func(w http.ResponseWriter, r *http.Request) {
		koalaId := chi.URLParam(r, "koalaId")

		depth := 2
		if d := r.URL.Query().Get("depth"); d != "" {
			// 雑に 0-3 のみ許可（深くしすぎると読みが増えるので）
			switch strings.TrimSpace(d) {
			case "0":
				depth = 0
			case "1":
				depth = 1
			case "2":
				depth = 2
			case "3":
				depth = 3
			}
		}

		includeChildren := true
		if v := strings.TrimSpace(r.URL.Query().Get("includeChildren")); v != "" {
			switch strings.ToLower(v) {
			case "0", "false", "no":
				includeChildren = false
			case "1", "true", "yes":
				includeChildren = true
			}
		}

		resp, err := buildPedigree(r.Context(), client, koalaId, depth, includeChildren)
		if err != nil {
			httpError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, resp)
	})

	r.Mount("/api", api)

	if cfg.ServeStatic {
		// static web (same origin)
		fileServer := http.FileServer(http.Dir(cfg.WebDir))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fileServer.ServeHTTP(w, r)
		})
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			applyStaticCacheHeaders(w, r.URL.Path)
			fileServer.ServeHTTP(w, r)
		})
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s (project=%s, serve_static=%t)", cfg.Port, cfg.ProjectID, cfg.ServeStatic)
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

func listKoalasByZoo(ctx context.Context, client *firestore.Client, zooId, status string) ([]Koala, error) {
	q := client.Collection("koalas").Where("currentZooId", "==", zooId)
	if status != "" {
		q = q.Where("status", "==", status)
	}

	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	out := make([]Koala, 0, len(docs))
	for _, d := range docs {
		m := d.Data()
		out = append(out, Koala{
			ID:           d.Ref.ID,
			Name:         asString(m["name"]),
			Sex:          asString(m["sex"]),
			BirthDate:    asString(m["birthDate"]),
			Status:       asString(m["status"]),
			CurrentZooID: asString(m["currentZooId"]),
			MotherID:     asString(m["motherId"]),
			FatherID:     asString(m["fatherId"]),
		})
	}
	return out, nil
}

func applyStaticCacheHeaders(w http.ResponseWriter, requestPath string) {
	ext := strings.ToLower(path.Ext(requestPath))
	if ext == "" {
		return
	}
	switch ext {
	case ".html":
		w.Header().Set("Cache-Control", "no-cache")
	default:
		w.Header().Set("Cache-Control", "public, max-age=86400")
	}
}

func buildPedigree(ctx context.Context, client *firestore.Client, rootID string, depth int, includeChildren bool) (PedigreeResponse, error) {
	// depth=0: 自分のみ
	// depth=1: 自分 + 両親
	// depth=2: 自分 + 両親 + 祖父母（親の両親）
	visited := map[string]bool{}
	nodesByID := map[string]PedigreeNode{}
	edges := make([]PedigreeEdge, 0)

	// 取得キュー（idと残り深さ）
	type item struct {
		id    string
		depth int
	}
	queue := []item{{id: rootID, depth: depth}}

	// まとめて取りたいので、必要なIDを集めてから都度バッチGetする方式
	// 小規模なので「足りないIDが出たら追加で取得」でもOK
	cache := map[string]map[string]any{} // id -> firestore data

	getDoc := func(id string) (map[string]any, bool, error) {
		if id == "" {
			return nil, false, nil
		}
		if v, ok := cache[id]; ok {
			return v, true, nil
		}
		ds, err := client.Collection("koalas").Doc(id).Get(ctx)
		if err != nil {
			// NotFound は「公開されてない/入ってない」扱いにする
			// FirestoreのNotFound判定は環境差があるので、ここでは文字列でゆるく判定
			if strings.Contains(strings.ToLower(err.Error()), "notfound") {
				return nil, false, nil
			}
			return nil, false, err
		}
		cache[id] = ds.Data()
		return cache[id], true, nil
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.id == "" || visited[cur.id] {
			continue
		}
		visited[cur.id] = true

		m, ok, err := getDoc(cur.id)
		if err != nil {
			return PedigreeResponse{}, err
		}
		if !ok {
			// データが無い（公開されてない等）→ unknownノードで置き換え
			nodesByID[cur.id] = PedigreeNode{
				ID:   cur.id,
				Name: "（不明）",
				Sex:  "U",
				Note: "record not found",
			}
			continue
		}

		// 自分ノード
		nodesByID[cur.id] = PedigreeNode{
			ID:        cur.id,
			Name:      asString(m["name"]),
			Sex:       asString(m["sex"]),
			BirthDate: asString(m["birthDate"]),
		}

		// 親を辿る
		if cur.depth > 0 {
			motherID := asString(m["motherId"])
			fatherID := asString(m["fatherId"])

			// mother edge
			if motherID != "" {
				edges = append(edges, PedigreeEdge{From: motherID, To: cur.id, Type: "mother"})
				queue = append(queue, item{id: motherID, depth: cur.depth - 1})
			} else {
				unk := "unknown:" + cur.id + ":mother"
				nodesByID[unk] = PedigreeNode{ID: unk, Name: "（母 不明）", Sex: "U"}
				edges = append(edges, PedigreeEdge{From: unk, To: cur.id, Type: "mother"})
			}

			// father edge
			if fatherID != "" {
				edges = append(edges, PedigreeEdge{From: fatherID, To: cur.id, Type: "father"})
				queue = append(queue, item{id: fatherID, depth: cur.depth - 1})
			} else {
				unk := "unknown:" + cur.id + ":father"
				nodesByID[unk] = PedigreeNode{ID: unk, Name: "（父 不明）", Sex: "U"}
				edges = append(edges, PedigreeEdge{From: unk, To: cur.id, Type: "father"})
			}
		}
	}

	// ルートの子どもも表示（対象コアラが親になっている場合）
	if includeChildren && rootID != "" {
		if err := addChildrenOfRoot(ctx, client, rootID, nodesByID, &edges); err != nil {
			return PedigreeResponse{}, err
		}
	}

	// nodesByID を配列化
	nodes := make([]PedigreeNode, 0, len(nodesByID))
	for _, n := range nodesByID {
		// 空の名前なら不明にしておく
		if strings.TrimSpace(n.Name) == "" {
			n.Name = "（不明）"
		}
		nodes = append(nodes, n)
	}

	return PedigreeResponse{
		RootID: rootID,
		Nodes:  nodes,
		Edges:  edges,
	}, nil
}

func addChildrenOfRoot(ctx context.Context, client *firestore.Client, rootID string, nodesByID map[string]PedigreeNode, edges *[]PedigreeEdge) error {
	if rootID == "" {
		return nil
	}

	addChild := func(d *firestore.DocumentSnapshot, relation string) {
		m := d.Data()
		id := d.Ref.ID
		if _, ok := nodesByID[id]; !ok {
			nodesByID[id] = PedigreeNode{
				ID:        id,
				Name:      asString(m["name"]),
				Sex:       asString(m["sex"]),
				BirthDate: asString(m["birthDate"]),
			}
		}
		*edges = append(*edges, PedigreeEdge{From: rootID, To: id, Type: relation})
	}

	motherDocs, err := client.Collection("koalas").Where("motherId", "==", rootID).Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	for _, d := range motherDocs {
		addChild(d, "mother")
	}

	fatherDocs, err := client.Collection("koalas").Where("fatherId", "==", rootID).Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	for _, d := range fatherDocs {
		addChild(d, "father")
	}
	return nil
}

func loadServerConfig() serverConfig {
	projectID := firstNonEmpty(
		strings.TrimSpace(os.Getenv("FIREBASE_PROJECT_ID")),
		strings.TrimSpace(os.Getenv("GOOGLE_CLOUD_PROJECT")),
		strings.TrimSpace(os.Getenv("GCP_PROJECT")),
	)
	if projectID == "" {
		log.Fatal("missing project id: set FIREBASE_PROJECT_ID or GOOGLE_CLOUD_PROJECT")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	serveStatic := parseBoolEnv("SERVE_STATIC", true)
	webDir := strings.TrimSpace(os.Getenv("WEB_DIR"))
	if webDir == "" {
		webDir = filepath.Join("..", "web")
	}

	if serveStatic {
		info, err := os.Stat(webDir)
		if err != nil {
			log.Fatalf("WEB_DIR does not exist: %s: %v", webDir, err)
		}
		if !info.IsDir() {
			log.Fatalf("WEB_DIR is not a directory: %s", webDir)
		}
	}

	return serverConfig{
		ProjectID:   projectID,
		Port:        port,
		ServeStatic: serveStatic,
		WebDir:      webDir,
	}
}

func parseBoolEnv(key string, defaultVal bool) bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if raw == "" {
		return defaultVal
	}
	switch raw {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		log.Fatalf("invalid boolean env %s=%q", key, raw)
		return defaultVal
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
