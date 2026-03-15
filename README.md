# コアラマップ JP

**日本全国のコアラ展示施設と血統をまとめたウェブサービス。**
動物園マップ・コアラ一覧・家系図ネットワークの3機能を提供する。

---

## 何ができるか

| 画面 | URL | 概要 |
|------|-----|------|
| 全国マップ | `/` | 日本地図上でコアラ在籍施設を確認。コアラ名で検索可能 |
| 動物園別一覧 | `/zoo.html?zooId=…` | 施設ごとのコアラアルバム。性別・状態でフィルタ |
| 個別家系図 | `/pedigree.html?koalaId=…` | 1頭を起点に親・子・パートナーをネットワーク表示 |
| 全体家系図 | `/pedigree-all.html` | 全コアラの血縁関係を一覧グラフで表示 |

---

## 技術構成

```
web/          静的 HTML（ビルドなし）
  index.html       全国マップ (Leaflet)
  zoo.html         動物園別一覧
  pedigree.html    個別家系図 (vis-network)
  pedigree-all.html 全体家系図 (vis-network)
  assets/          アイコン・画像

api/          Go API サーバー
  cmd/server/      HTTP サーバー (chi v5)
  cmd/seed/        Firestore 初期データ投入

infra/        デプロイスクリプト
  deploy-prod.sh   Cloud Run + Hosting 一括デプロイ
```

**スタック:** Go 1.25 / chi v5 / Cloud Firestore / Firebase Hosting / Cloud Run
**フロント:** Vanilla JS + Leaflet + vis-network（ビルドステップなし）

### API エンドポイント

| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/health` | ヘルスチェック |
| GET | `/api/zoos` | 施設一覧 |
| GET | `/api/zoos/{zooId}/koalas` | 施設別コアラ一覧 (`?status=alive`) |
| GET | `/api/koalas/{koalaId}/pedigree` | 家系図 (`?depth=0..3&includeChildren=0\|1`) |
| GET | `/api/pedigree/all` | 全コアラ血縁グラフ |

### データモデル

**`zoos`** — `name`, `pref`, `city`, `status`, `officialUrl`, `lastVerifiedAt`, `geo{lat,lng}`, `sources[]`
**`koalas`** — `name`, `sex`, `birthDate`, `status`, `currentZooId`, `motherId`, `fatherId`, `lastVerifiedAt`, `photoUrl`, `sources[]`

---

## ローカル開発

### A. Docker（推奨）

```bash
# 全サービス起動（Firestore Emulator + API + Nginx）
docker compose up --build -d firestore api web

# 初回: シードデータ投入
docker compose run --rm seed
```

| URL | 内容 |
|-----|------|
| `http://localhost:8088` | 画面 |
| `http://localhost:8088/api/health` | API ヘルスチェック |
| `http://localhost:4000` | Firestore Emulator UI |

```bash
docker compose down   # 停止（データは .emulator-data/ に保存）
```

### B. ローカル直接実行

**前提:** Go / Firebase CLI インストール済み

```bash
# 1. Firestore Emulator 起動
firebase emulators:start --only firestore --import=.emulator-data --export-on-exit

# 2. API サーバー起動（別ターミナル）
export FIREBASE_PROJECT_ID=koala-map-jp
export FIRESTORE_EMULATOR_HOST=127.0.0.1:8081
cd api && go run ./cmd/server

# 3. 初回: シードデータ投入（別ターミナル）
export FIREBASE_PROJECT_ID=koala-map-jp
export FIRESTORE_EMULATOR_HOST=127.0.0.1:8081
cd api && go run ./cmd/seed
```

---

## 本番デプロイ（Cloud Run + Firebase Hosting）

**前提:** `gcloud auth login` / `firebase login` 済み

```bash
./infra/deploy-prod.sh
# または
PROJECT_ID=koala-map-jp REGION=asia-northeast1 API_SERVICE=koala-map-jp-api ./infra/deploy-prod.sh
```

疎通確認:

```bash
curl -i https://koala-map-jp.web.app/api/health
```

**構成:** Firebase Hosting が `/api/**` を Cloud Run (`koala-map-jp-api`) へ rewrite。
Firestore ルールはクライアント SDK からの直接アクセスを拒否（API 経由のみ）。

---

## 環境変数

| 変数 | デフォルト | 説明 |
|------|-----------|------|
| `FIREBASE_PROJECT_ID` | `GOOGLE_CLOUD_PROJECT` の値 | Firestore プロジェクト ID |
| `PORT` | `8080` | API リッスンポート |
| `SERVE_STATIC` | `true` | `true` で web/ を同梱配信 |
| `WEB_DIR` | `../web` | 静的ファイルのパス |
| `FIRESTORE_EMULATOR_HOST` | （未設定） | Emulator 使用時に設定 |
