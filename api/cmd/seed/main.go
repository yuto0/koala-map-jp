package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

const dateLayout = "2006-01-02"

func main() {
	ctx := context.Background()

	projectID := firstNonEmpty(
		strings.TrimSpace(os.Getenv("FIREBASE_PROJECT_ID")),
		strings.TrimSpace(os.Getenv("GOOGLE_CLOUD_PROJECT")),
		strings.TrimSpace(os.Getenv("GCP_PROJECT")),
	)
	if projectID == "" {
		log.Fatal("missing project id: set FIREBASE_PROJECT_ID or GOOGLE_CLOUD_PROJECT")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	today := time.Now().Format(dateLayout)
	zoos := buildZoos(today)
	koalas := buildKoalas(today)

	// -----------------------
	// Firestoreへ投入（MergeAll）
	// -----------------------
	batch := client.Batch()

	// zoos
	for _, z := range zoos {
		ref := client.Collection("zoos").Doc(z.ID)
		batch.Set(ref, z.toDoc(today), firestore.MergeAll)
	}

	// koalas
	for _, k := range koalas {
		ref := client.Collection("koalas").Doc(k.ID)
		batch.Set(ref, k.toDoc(today), firestore.MergeAll)
	}

	if _, err := batch.Commit(ctx); err != nil {
		log.Fatalf("batch.Commit: %v", err)
	}

	log.Printf("seeded zoos=%d, koalas=%d", len(zoos), len(koalas))
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
