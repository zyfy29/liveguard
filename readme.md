# Liveguard

Liveguard は、セキュリティ監視システムのバックエンドとフロントエンドからなるプロジェクトです。

## 構成

- **Backend (back)**
    - データ処理・API 通信を担当
    - `config.yaml`（`config.example.yaml` を元に作成）に依存
    - 必要に応じて SQLite3 の DB ファイル `liveguard.db` を利用

- **Frontend (front)**
    - Streamlit を用いた Web UI で、バックエンドからのデータを可視化
    - 環境変数 `DYNACONF_API_ADDR` と `DYNACONF_API_PORT` を利用してバックエンドに接続

## デプロイ時のディレクトリ構成

```
liveguard-deploy/
├── config.yaml       # デプロイ前に作成
├── liveguard.db      # 必要なら作成 (touch liveguard.db)
├── .env              # ユーザーがポートを指定するためのファイル
└── docker-compose.yaml
```

## Docker Compose を使ったデプロイ方法

### 前提条件

- Docker および Docker Compose がインストール済みであること
- `back/config.example.yaml` を元に `config.yaml` を作成する
- 必要に応じて、SQLite3 用の DB ファイルを作成する

  ```sh
  touch liveguard.db
  ```
- `back/config.example.yaml` を元に `.env` を作成する
- `.env` 内の `BACK_PORT` と `FRONT_PORT` を適切に設定する

### デプロイ手順

1. **イメージのビルド**

   ```sh
   docker compose build
   ```

2. **サービスの起動**

   ```sh
   docker compose up -d
   ```

3. **アクセス方法**

    - **Backend:** `BACK_PORT` で起動
    - **Frontend:** `FRONT_PORT` で公開され、Web ブラウザからアクセス可能
