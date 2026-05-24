# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## コマンド

タスクランナーは [Task](https://taskfile.dev) (`go-task`)。`Taskfile.yml` を参照。

- **サーバ起動**: `task run` (`:8080` でリッスン。事前に `task mongo:up` 必要)
- **バイナリビルド**: `task build`
- **テスト実行**: `task test` (`testcontainers-go` が `mongo:6.0.2` を自動起動)
- **静的解析**: `task vet`
- **モジュール整理**: `task tidy`
- **開発用 MongoDB 起動 / 停止**: `task mongo:up` / `task mongo:down`
- **タスク一覧**: `task`

`Taskfile.yml` は `docker context inspect | jq` で現在の docker context から `DOCKER_HOST` を自動取得しているので、Docker Desktop / colima / CI で共通の Taskfile が動く。`TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock` も全タスクで export 済み。

特定のテストだけ走らせる場合は task を経由せず直接 `go test ./handler -run TestRegisterUser` のように呼べばよい (env は `task test` と同等を別途整える必要あり)。

## アーキテクチャ

Go 標準ライブラリ (`net/http`) で構成された単一バイナリの HTTP サービス。Web フレームワークやルータライブラリは使っていない。ルーティングは `main.go` 内で `http.HandleFunc` を使って直接登録している。

- `main.go` — ルート登録とサーバ起動。`/user` ルートはローカルの `userHandler` 内で HTTP メソッドごとに switch で振り分ける構造になっており、メソッドごとに別パスを登録するスタイルではない点に注意。
- `handler/` パッケージ — リクエストハンドラ一式を配置。
  - `base.go` — MongoDB 接続に関する設定 (`uri`, `databaseName`, `userCollection`) と JSON レスポンスヘルパ (`response`, `responseListUsers`)。`uri` は `var` で宣言されており、テスト (`main_test.go`) が testcontainers の接続文字列で上書きする想定。
  - `sample.go` — `Index` (`/`、`?name=` 必須) と `Ping` (`/db/ping`、Mongo 接続確認用)。
  - `user.go` — `ListUser` (GET `/users`) と `RegisterUser` (POST `/user`)、および `User` 構造体 (`Id` は `bson.ObjectID`)。

### MongoDB の利用パターン

各ハンドラがそれぞれ `mongo.Connect` を呼んで `Disconnect` を defer しており、共有クライアントやコネクションプールはハンドラ間で共有していない。Mongo を叩く新しいハンドラを追加するときは `ListUser`/`RegisterUser` と同じ connect-defer-disconnect パターンを踏襲するか、性能が問題になるなら共有クライアントとして切り出すこと。

`User.Id` は `bson.ObjectID` 型で、`RegisterUser` は挿入時に値をセットしない。Mongo 側で `ObjectID` が自動生成され、ハンドラはその hex 文字列をレスポンスメッセージとして返す。

### テスト構成 (handler パッケージ)

- `main_test.go` の `TestMain` が `mongo:6.0.2` コンテナを testcontainers-go で 1 度だけ起動し、`uri` を上書きしてから `m.Run()` を呼ぶ。終了時に terminate。
- DB 依存テスト (`user_test.go`, `ping_test.go`) は `resetUserCollection(t)` ヘルパで pre/post に `user` コレクションを drop し、テスト間の独立性を確保。
- DB 非依存の `Index` テスト (`sample_test.go`) は `httptest.ResponseRecorder` のみ。
