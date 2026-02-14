# koala-map-jp 現状整理

最終更新: 2026-02-14

## 1. アプリ構成
- `api/`: Go 製 API サーバー
- `web/`: 静的フロントエンド（`index.html`, `zoo.html`, `pedigree.html`）
- `firestore.rules` / `firestore.indexes.json`: Firestore 設定
- `firebase.json`: Hosting / Emulator 設定

## 2. 実行エントリポイント
- API サーバー: `api/cmd/server/main.go`
- Seed 投入: `api/cmd/seed/main.go`

必要環境変数:
- `FIREBASE_PROJECT_ID`（未指定時は `GOOGLE_CLOUD_PROJECT` を使用）
- `PORT`（未指定時 `8080`）
- `SERVE_STATIC`（未指定時 `true`）
- `WEB_DIR`（未指定時 `../web`）

## 3. API 一覧（現状）
- `GET /api/health`
- `GET /api/zoos`
- `GET /api/zoos/{zooId}/koalas?status=alive`
- `GET /api/koalas/{koalaId}/pedigree?depth=0..3&includeChildren=0|1`

## 4. フロント画面（現状）
- `web/index.html`: 全国施設マップ + 一覧
- `web/zoo.html`: 動物園別コアラ一覧
- `web/pedigree.html`: 家系図ネットワーク表示

フロントは同一オリジンで `/api` を呼ぶ構成。

## 5. Firestore データモデル（現状）
- コレクション `zoos`
  - 主なフィールド: `name`, `pref`, `city`, `status`, `officialUrl`, `lastVerifiedAt`, `geo{lat,lng}`, `sources[]`
- コレクション `koalas`
  - 主なフィールド: `name`, `sex`, `birthDate`, `status`, `currentZooId`, `motherId`, `fatherId`, `lastVerifiedAt`, `photoUrl`, `sources[]`

## 6. リポジトリ状態（2026-02-14 時点）
未コミット変更あり:
- 変更: `api/cmd/seed/main.go`, `api/cmd/server/main.go`, `web/index.html`
- 新規: `.vscode/`, `api/cmd/seed/data.go`, `image.png`, `kakeizu.png`, `web/assets/`, `web/pedigree.html`, `web/zoo.html`

## 7. 本番向け対応の反映状況
1. API 実行モード分離
   - `SERVE_STATIC=true/false` で、ローカル一体運用と Cloud Run API 専用運用を切替可能。
2. Hosting から API への経路固定
   - `firebase.json` の rewrite で `/api/**` を Cloud Run (`koala-map-jp-api`) に転送。
3. Firestore ルール最小化
   - `firestore.rules` はクライアントSDKからの read/write を拒否し、API 経由アクセス前提。
4. デプロイ自動化
   - `infra/deploy-prod.sh` で Cloud Run / Firestore / Hosting を順次デプロイ可能。

## 8. 技術的リスク / 要整理ポイント
1. ルート直下にログファイルが残存
   - `firebase-debug.log`, `firestore-debug.log` が大きく、ローカル確認ノイズになる。
2. 依存管理の可読性
   - `api/go.mod` は全依存が `indirect` 表示で、主要直接依存が分かりにくい。
3. 運用監視が未整備
   - Cloud Run / Firestore のエラー監視、予算アラート、バックアップ運用が未定義。

## 9. 直近の優先アクション
1. 本番デプロイ実行（`./infra/deploy-prod.sh`）と疎通確認（`/api/health`）。
2. Cloud Run サービスアカウント権限を最小化（Firestore必要権限のみに限定）。
3. Cloud Logging のアラート（5xx / エラーレート）を設定。
4. Firestore の定期バックアップ運用を決定。
