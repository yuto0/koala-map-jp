# koala-map-jp

全国のコアラ展示施設を地図で見られるサイト（学習用に Go + Firestore で構築）。

## 現状整理
- 現状サマリ: `CURRENT_STATE.md`

## 構成
- `api/`: Go API
- `web/`: 静的フロントエンド
- `firestore.rules` / `firestore.indexes.json`: Firestore 設定
- `firebase.json`: Hosting / Emulator 設定

## ローカル実行（最小）
前提:
- Go が利用可能
- `FIREBASE_PROJECT_ID` を設定済み

API 起動:
```bash
cd api
FIREBASE_PROJECT_ID=koala-map-jp go run ./cmd/server
```

Seed 投入:
```bash
cd api
FIREBASE_PROJECT_ID=koala-map-jp go run ./cmd/seed
```

## ローカルDB永続化（Firestore Emulator）
前提:
- Firebase CLI が利用可能

Firestore Emulator をデータ永続化付きで起動:
```bash
firebase emulators:start --only firestore --import=.emulator-data --export-on-exit
```

動作:
- 初回終了時に `.emulator-data/` が作成される
- 次回起動時に `.emulator-data/` から自動復元される

API を Emulator に接続して起動:
```bash
export FIREBASE_PROJECT_ID=koala-map-jp
export FIRESTORE_EMULATOR_HOST=127.0.0.1:8081
cd api
go run ./cmd/server
```

Seed を Emulator に投入:
```bash
export FIREBASE_PROJECT_ID=koala-map-jp
export FIRESTORE_EMULATOR_HOST=127.0.0.1:8081
cd api
go run ./cmd/seed
```

## 本番デプロイ（Cloud Run + Hosting）
前提:
- `gcloud` / `firebase` CLI が利用可能
- `gcloud auth login` 済み
- `firebase login` 済み

現在の `firebase.json` は `/api/**` を Cloud Run の `koala-map-jp-api` に rewrite する設定です。

一括デプロイ:
```bash
./infra/deploy-prod.sh
```

環境を上書きしてデプロイ:
```bash
PROJECT_ID=koala-map-jp REGION=asia-northeast1 API_SERVICE=koala-map-jp-api ./infra/deploy-prod.sh
```

デプロイ後の確認:
```bash
curl -i https://koala-map-jp.web.app/api/health
```
