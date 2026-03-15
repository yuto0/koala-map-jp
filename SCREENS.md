# 画面設計書

## 画面一覧

| 画面名 | ファイル | URL | 概要 |
|--------|---------|-----|------|
| トップ（地図） | `web/index.html` | `/` | 全国コアラ地図・検索 |
| 動物園詳細 | `web/zoo.html` | `/zoo.html?zooId={id}` | 動物園のコアラ一覧 |
| 個別家系図 | `web/pedigree.html` | `/pedigree.html?koalaId={id}` | 個体の家系図グラフ |
| 全体家系図 | `web/pedigree-all.html` | `/pedigree-all.html` | 全コアラ統合グラフ |
| 404 | `web/404.html` | — | Not Found ページ |

---

## 画面遷移

```
index.html
  └─ 動物園カードクリック → zoo.html?zooId=...
       └─ コアラカードクリック → pedigree.html?koalaId=...&zooId=...
            ├─ ノードクリック → pedigree.html?koalaId=... （別個体）
            └─ 全体図へ → pedigree-all.html
                 └─ ノードクリック → pedigree.html?koalaId=...
```

---

## 各画面詳細

### index.html — トップ（地図）

**目的**: 日本国内のコアラ所在地を地図と一覧で把握する

**使用ライブラリ**: Leaflet.js

**API**:
- `GET /api/zoos` — 全動物園取得（地図マーカー・一覧表示）
- `GET /api/zoos/{zooId}/koalas` — 動物園ポップアップ内のコアラ一覧

**主な機能**:
- 地図マーカー（コアラ数に応じてサイズ変動）
- コアラ名によるリアルタイム検索
- 動物園一覧の並び順切り替え（あいうえお順 / コアラ数順）
- ヘッダー統計（動物園数、登録コアラ数、最終確認日）
- マーカークリックでポップアップ（前後の動物園ナビゲーション付き）

**遷移先**: `zoo.html?zooId={id}`

---

### zoo.html — 動物園詳細

**目的**: 指定動物園のコアラ一覧を検索・絞り込みして確認する

**URLパラメータ**:
| パラメータ | デフォルト | 説明 |
|-----------|-----------|------|
| `zooId` | `higashiyama-zoo` | 動物園ID |

**API**:
- `GET /api/zoos` — 動物園情報（名前・公式サイトURL等）
- `GET /api/zoos/{zooId}/koalas` — コアラ一覧

**主な機能**:
- 名前検索
- 性別フィルター（すべて / オス / メス / 不明）
- 状態フィルター（すべて / 現存 / 死亡 / 不明）
- 並び順（名前 / 生年新しい順 / 生年古い順）
- 統計バッジ（現存数、メス数、オス数）
- コアラカードのモーダル詳細表示

**遷移先**: `pedigree.html?koalaId={id}&depth=2&includeChildren=1&zooId={id}`

---

### pedigree.html — 個別家系図

**目的**: 特定コアラの祖先・子孫関係をグラフで可視化する

**URLパラメータ**:
| パラメータ | デフォルト | 説明 |
|-----------|-----------|------|
| `koalaId` | `monaka` | 中心となるコアラID |
| `depth` | `2` | 遡る世代数（0〜3） |
| `includeChildren` | `1` | 子世代を含める（0/1） |
| `zooId` | — | 戻り先の動物園ID（任意） |

**API**:
- `GET /api/koalas/{koalaId}/pedigree?depth={n}&includeChildren={0|1}`

**レスポンス**:
```json
{
  "rootId": "string",
  "nodes": [{ "id", "name", "sex", "birthDate", "status", "hasParents" }],
  "edges": [{ "from", "to", "type": "mother|father" }]
}
```

**グラフ描画**:
- vis-network（hierarchical / UD方向 / physics無効）
- 両親揃いの場合: 母→ユニオンノード←父、ユニオンノード→子
- 片親のみの場合: 親→子 直接エッジ
- ノード色: メス=#b84c7a、オス=#2e6ab0、ルート=#8b2a1b

**インタラクション**:
- ノードクリック → 別コアラの家系図へ遷移
- 性別・生年月日の表示切り替え
- 中心へ / 全体表示ボタン

**遷移先**:
- `zoo.html?zooId={id}` または `/` （戻る）
- `pedigree-all.html` （全体図へ）
- `pedigree.html?koalaId={id}` （別個体）

---

### pedigree-all.html — 全体家系図

**目的**: DB内の全コアラの血統関係を一画面で俯瞰する

**URLパラメータ**:
| パラメータ | デフォルト | 説明 |
|-----------|-----------|------|
| `zooId` | — | 戻り先用（任意） |

**API**:
- `GET /api/pedigree/all`

**グラフ描画**:
- vis-network（physics有効 / force-directed）
- pedigree.html と同じユニオンノード構造
- 安定化完了後に自動フィット

**インタラクション**:
- 名前でハイライト検索（ノードの強調表示）
- 性別・生年月日の表示切り替え
- 全体表示ボタン
- ノードクリック → 個別家系図へ遷移

**遷移先**: `pedigree.html?koalaId={id}&depth=2&includeChildren=1`

---

### 404.html — Not Found

**目的**: 存在しないURLへのアクセスをホームへ誘導する

静的ページ。APIなし。ホームへのリンクのみ。

---

## APIエンドポイント一覧

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/health` | ヘルスチェック |
| GET | `/api/zoos` | 全動物園一覧 |
| GET | `/api/zoos/{zooId}/koalas` | 動物園のコアラ一覧（`?status=` 任意） |
| GET | `/api/koalas/{koalaId}/pedigree` | 個体の家系図（`?depth=&includeChildren=`） |
| GET | `/api/pedigree/all` | 全体家系図 |

**定義ファイル**: `api/cmd/server/main.go`
