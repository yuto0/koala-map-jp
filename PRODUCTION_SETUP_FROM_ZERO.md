# Production Setup From Zero

Firebase プロジェクトは作成済みで、GCP 側の本番実行基盤（Cloud Run / Firestore / IAM / 課金）が未設定の状態から始める手順。

対象リポジトリ: `koala-map-jp`

## 0. 前提
- 本番プロジェクトID: `koala-map-jp`
- テストプロジェクトID: `koala-map-jp-test`（推奨）
- 利用リージョン: `asia-northeast1`
- 本番APIサービス名: `koala-map-jp-api`
- テストAPIサービス名: `koala-map-jp-api-test`
- ローカルに `gcloud` / `firebase` が入っている

## 1. 認証とプロジェクト選択
```bash
gcloud auth login
firebase login

gcloud config set project koala-map-jp
firebase use koala-map-jp
```

確認:
```bash
gcloud config get-value project
firebase projects:list | rg koala-map-jp
```

## 2. 課金アカウントを Firebase/GCP プロジェクトに紐付け
Cloud Run を使うには Blaze（従量課金）が必要。

1. Firebase Console で `koala-map-jp` を Blaze にアップグレード
2. GCP Console 側で Billing アカウントが `koala-map-jp` にリンクされていることを確認

CLI確認:
```bash
gcloud beta billing projects describe koala-map-jp
```

## 3. 必要 API を有効化
```bash
gcloud services enable \
  run.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  firestore.googleapis.com \
  iam.googleapis.com \
  serviceusage.googleapis.com
```

## 4. Firestore データベース作成（未作成の場合のみ）
Firestore が未作成なら実行。

```bash
gcloud firestore databases create \
  --location=asia-northeast1 \
  --type=firestore-native
```

作成済み確認:
```bash
gcloud firestore databases list --project koala-map-jp
```

## 5. サービスアカウント権限の確認（最小）
まずは現状動作を優先して Cloud Run 実行サービスアカウントを確認。

```bash
gcloud run services list --region asia-northeast1
```

初回デプロイ後、実行SA（通常は Compute Default か指定SA）に対し Firestore 読み取りに必要なロールを付与する。

例（必要最小の一例）:
```bash
gcloud projects add-iam-policy-binding koala-map-jp \
  --member="serviceAccount:YOUR_RUNTIME_SA" \
  --role="roles/datastore.user"
```

## 6. テスト環境へデプロイ（先行）
リポジトリルートで実行:
```bash
cd /Users/yuto/Developer/koala-map-jp
./infra/deploy-test.sh
```

環境を明示する場合:
```bash
PROJECT_ID=koala-map-jp-test REGION=asia-northeast1 API_SERVICE=koala-map-jp-api-test ./infra/deploy-test.sh
```

テスト疎通確認:
```bash
curl -i https://koala-map-jp-test.web.app/api/health
```

期待値:
- `HTTP/2 200`
- body: `ok`

## 7. 本番デプロイ（テスト確認後）
リポジトリルートで実行:
```bash
cd /Users/yuto/Developer/koala-map-jp
./infra/deploy-prod.sh
```

環境を明示する場合:
```bash
PROJECT_ID=koala-map-jp REGION=asia-northeast1 API_SERVICE=koala-map-jp-api ./infra/deploy-prod.sh
```

このスクリプトは以下を順に実行する:
1. Cloud Run デプロイ
2. Firestore rules/indexes デプロイ
3. Hosting デプロイ

## 8. 本番デプロイ後の疎通確認
```bash
curl -i https://koala-map-jp.web.app/api/health
```

期待値:
- `HTTP/2 200`
- body: `ok`

追加確認:
```bash
# Hostingトップ
curl -I https://koala-map-jp.web.app/

# Cloud Runサービス一覧
gcloud run services list --region asia-northeast1 --project koala-map-jp
```

## 9. 失敗時の切り分け
### A. `PERMISSION_DENIED`
- 実行SAの IAM ロール不足を確認
- `roles/datastore.user` が付いているか確認

### B. `Billing account not configured`
- Blaze へのアップグレード漏れ
- Billing リンク未設定

### C. `API not enabled`
- 手順3の `gcloud services enable` を再実行

### D. `/api/health` が 404/502
- `firebase.json` の `/api/**` rewrite と Cloud Run サービス名/リージョンを確認

## 10. 最初にやっておく運用設定（推奨）
1. 予算アラート作成（少額閾値）
2. Cloud Run 5xx アラート作成
3. Firestore 使用量アラート作成
4. 本番用 Runtime Service Account を専用化

---
この手順は「初回セットアップ」に特化。日常の更新は `DEVELOPMENT_WORKFLOW.md` と `./infra/deploy-test.sh` / `./infra/deploy-prod.sh` を利用する。
