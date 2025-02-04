# Liveguard Panel

Liveguard Panel は、セキュリティ監視システムのフロントエンドです。

## 🚀 技術スタック
- **Python 3.12**: 最新の Python を使用
- **Streamlit**: 軽量でシンプルな Web UI フレームワーク
- **Poetry**: 依存管理とパッケージ管理
- **API クライアント**: `api_client.py` でバックエンドとの通信を担当

## 🎯 Streamlit を選んだ理由
1. **素早く UI を作成できる**
    - HTML/CSS/JS の知識がなくても、Python だけで Web UI を構築可能
2. **データ可視化が簡単**
    - グラフやテーブル表示が標準でサポートされている
3. **開発が速い**
    - コードを少なくできるため、PoC（概念実証）や MVP 開発に最適
4. **API クライアントとの統合が容易**
    - `requests` などを使って API からデータを取得し、リアルタイムで表示可能
