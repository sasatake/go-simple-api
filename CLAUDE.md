# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## コマンド

- **サーバ起動**: `go run main.go` (`:8080` でリッスン)
- **バイナリビルド**: `go build -o go-simple-api`
- **MongoDB 起動**: `docker compose up -d` (mongo 6.0.2 が `localhost:27017`、認証は `mongo`/`mongo`)
- **MongoDB 停止**: `docker compose down`
- **モジュール整理**: `go mod tidy`
- **静的解析**: `go vet ./...`

このリポジトリにはまだテストスイートや lint 設定はありません。

## アーキテクチャ

Go 標準ライブラリ (`net/http`) で構成された単一バイナリの HTTP サービスで、Web フレームワークやルータライブラリは使っていません。ルーティングは `main.go` 内で `http.HandleFunc` を使って直接登録しています。

- `main.go` — ルート登録とサーバ起動。`/user` ルートはローカルの `userHandler` 内で HTTP メソッドごとに switch で振り分ける構造になっており、メソッドごとに別パスを登録するスタイルではない点に注意。
- `handler/` パッケージ — リクエストハンドラ一式を配置。
  - `base.go` — MongoDB 接続に関する共通定数 (`uri`, `databaseName`, `userCollection`) と JSON レスポンスヘルパ (`response`, `responseListUsers`)。Mongo の接続文字列はここにハードコードされているので、接続先を変える場合はここを修正する。
  - `sample.go` — `Index` (`/`、`?name=` 必須) と `Ping` (`/db/ping`、Mongo 接続確認用)。
  - `user.go` — `ListUser` (GET `/users`) と `RegisterUser` (POST `/user`)、および `json` と `bson` のタグを併用した `User` 構造体。

### MongoDB の利用パターン

各ハンドラがそれぞれ `mongo.Connect` を呼んで `Disconnect` を defer しており、共有クライアントやコネクションプールはハンドラ間で共有していません。Mongo を叩く新しいハンドラを追加するときは `ListUser`/`RegisterUser` と同じ connect-defer-disconnect パターンを踏襲するか、性能が問題になるなら共有クライアントとして切り出してください。

`User.Id` フィールドには `bson:"_id"` タグが付いていますが、`RegisterUser` は挿入時に Id をセットしません。Mongo 側で `ObjectID` が自動生成され、ハンドラはその hex 文字列をレスポンスメッセージとして返します。
