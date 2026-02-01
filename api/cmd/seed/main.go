package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
)

type Zoo struct {
	ID             string
	Name           string
	Pref           string
	City           string
	Status         string
	OfficialURL    string
	LastVerifiedAt string
	Lat            float64
	Lng            float64
}

func main() {
	ctx := context.Background()

	projectID := mustEnv("FIREBASE_PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	// 今日は固定で入れてOK（必要なら後で更新する運用に）
	today := time.Now().Format("2006-01-02")

	zoos := []Zoo{
		{
			ID: "hirakawa-zoo",
			Name: "平川動物園",
			Pref: "鹿児島県", City: "鹿児島市",
			Status: "open",
			OfficialURL: "https://hirakawazoo.jp/",
			LastVerifiedAt: today,
			Lat: 31.463123, Lng: 130.503422,
		},
		{
			ID: "england-hill",
			Name: "淡路ファームパーク イングランドの丘",
			Pref: "兵庫県", City: "南あわじ市",
			Status: "open",
			OfficialURL: "https://www.england-hill.com/",
			LastVerifiedAt: today,
			Lat: 34.3071605, Lng: 134.8017872,
		},
		{
			ID: "oji-zoo",
			Name: "王子動物園",
			Pref: "兵庫県", City: "神戸市",
			Status: "open",
			OfficialURL: "https://www.kobe-ojizoo.jp/",
			LastVerifiedAt: today,
			Lat: 34.7099614, Lng: 135.2145663,
		},
		{
			ID: "tama-zoo",
			Name: "多摩動物園",
			Pref: "東京都", City: "日野市",
			Status: "open",
			OfficialURL: "https://www.tokyo-zoo.net/zoo/tama/",
			LastVerifiedAt: today,
			Lat: 35.649538, Lng: 139.402201,
		},
		{
			ID: "kanazawa-zoo",
			Name: "金沢動物園",
			Pref: "神奈川県", City: "横浜市",
			Status: "open",
			OfficialURL: "https://www.hama-midorinokyokai.or.jp/zoo/kanazawa/",
			LastVerifiedAt: today,
			Lat: 35.349729, Lng: 139.597359,
		},
		{
			ID: "higashiyama-zoo",
			Name: "東山動植物園",
			Pref: "愛知県", City: "名古屋市",
			Status: "open",
			OfficialURL: "https://www.higashiyama.city.nagoya.jp/",
			LastVerifiedAt: today,
			Lat: 35.1567406, Lng: 136.9811452,
		},
		{
			ID: "saitama-childrens-zoo",
			Name: "さいたまこども動物園（埼玉県こども動物自然公園）",
			Pref: "埼玉県", City: "東松山市",
			Status: "open",
			OfficialURL: "https://www.parks.or.jp/sczoo/",
			LastVerifiedAt: today,
			Lat: 36.000865, Lng: 139.374657,
		},
	}

	batch := client.Batch()
	for _, z := range zoos {
		ref := client.Collection("zoos").Doc(z.ID)
		batch.Set(ref, map[string]any{
			"name":           z.Name,
			"pref":           z.Pref,
			"city":           z.City,
			"status":         z.Status,
			"officialUrl":    z.OfficialURL,
			"lastVerifiedAt": z.LastVerifiedAt,
			"geo": map[string]any{
				"lat": z.Lat,
				"lng": z.Lng,
			},
		})
	}

	if _, err := batch.Commit(ctx); err != nil {
		log.Fatalf("batch.Commit: %v", err)
	}

	log.Printf("seeded %d zoos into collection 'zoos'", len(zoos))
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env: %s", key)
	}
	return v
}
