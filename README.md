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

## Docker 擬似本番環境（Hosting + API + Firestore）
この構成は `web(nginx)` -> `/api` プロキシ -> `Go API` -> `Firestore Emulator` です。
Firestore Emulator コンテナは `Java 21` を利用します。

起動:
```bash
docker compose up --build -d firestore api web
```

Seed投入:
```bash
docker compose run --rm seed
```

アクセス先:
- 画面: `http://localhost:8088`
- APIヘルス: `http://localhost:8088/api/health`
- Firestore Emulator UI: `http://localhost:4000`

停止:
```bash
docker compose down
```

停止時に Firestore データは `.emulator-data/` へエクスポートされ、次回起動時に復元されます。

Firestore コンテナ設定を変更したときは再ビルドしてください:
```bash
docker compose build firestore
docker compose up -d firestore
```
