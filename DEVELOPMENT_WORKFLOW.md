# Development Workflow

このプロジェクトで効率的に開発するための標準手順。

## 1. 作業開始
1. `main` を最新化する。
2. 作業ブランチを作成する。

## 2. ローカル基盤起動
Firestore Emulator:
```bash
firebase emulators:start --only firestore --import=.emulator-data --export-on-exit
```

## 3. API起動（別ターミナル）
```bash
export FIREBASE_PROJECT_ID=koala-map-jp
export FIRESTORE_EMULATOR_HOST=127.0.0.1:8081
cd api
go run ./cmd/server
```

## 4. Seed投入（必要時）
```bash
cd api
FIREBASE_PROJECT_ID=koala-map-jp FIRESTORE_EMULATOR_HOST=127.0.0.1:8081 go run ./cmd/seed
```

## 5. 実装ループ
1. 小さく実装する。
2. API疎通を確認する。
3. 変更を見直す。

ヘルスチェック:
```bash
curl -i http://localhost:8080/api/health
```

## 6. 品質ゲート
```bash
cd api
go test ./...
```

必要時のみ整形:
```bash
cd api
gofmt -w ./cmd
```

## 7. 仕上げ
1. 変更内容を要約する。
2. 確認手順と未解決リスクを整理する。
3. コミットする。

## 8. デプロイ時
```bash
./infra/deploy-prod.sh
curl -i https://koala-map-jp.web.app/api/health
```
