---
name: git-commit
description: このリポジトリの変更を日本語コミットメッセージでコミットし、続けて現在ブランチを push する。status / diff / log を確認したうえで、変更内容に沿った簡潔な日本語メッセージで新しいコミットを作り、最後に origin へ push する。ユーザが「コミットして」「commit して」「push して」などコミットや push を依頼したときに使う。
---

# git-commit

このリポジトリの変更内容を確認し、日本語のコミットメッセージで新しいコミットを作成し、続けて push まで行うためのスキル。

## 手順

以下の手順を順守すること。

### 1. 現状把握 (並列実行)

次の 3 つを **1 つのメッセージで並列に** Bash 実行する:

- `git status` — 変更ファイル一覧 (`-uall` は使わない)
- `git diff` — staged / unstaged の両方を含む差分
- `git log --oneline -10` — このリポジトリのコミットメッセージのスタイルを把握する

このリポジトリの既存コミット例 (`git log` で確認):

```
de15ce8 create handler ListUser and RegisterUser
a094afb add mongodb connection handler
e145abd response to json
6b9553d add simple handler
70c28f0 go mod init
```

→ 既存履歴は英語の短い動詞始まりだが、**本スキルでは明示の指示に従い日本語で書く**。スタイルは「動詞 + 対象」を意識した短い 1 行を基本とする。

### 2. メッセージのドラフト

差分から以下を判定してメッセージを組み立てる:

- 変更の種類を表す動詞を冒頭に置く:
  - 新機能: `追加:` (例: `追加: /user の DELETE ハンドラ`)
  - 既存機能の拡張・改善: `更新:`
  - バグ修正: `修正:`
  - リファクタ: `リファクタ:`
  - ドキュメント: `docs:`
  - テスト: `test:`
  - 設定・ビルド系: `chore:`
- 件名は 50 文字程度に収める。長い説明は本文 (空行を 1 行挟んだあと) に書く。
- 「何を」だけでなく「なぜ」が非自明なら本文に書く。自明なら件名のみで OK。
- `.env` などシークレットを含みうるファイルは原則ステージしない。ユーザが明示的に指示した場合のみ警告したうえで含める。

### 3. ステージング + コミット (並列実行)

- `git add <files>` — `git add -A` や `git add .` は使わず、対象ファイルを明示する
- `git commit -m "..."` — 複数行メッセージは必ず HEREDOC で渡す:

```bash
git commit -m "$(cat <<'EOF'
追加: ユーザ削除ハンドラ

DELETE /user に対応。ObjectID 指定で 1 件削除する。

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

- コミット完了後に `git status` を実行して結果を確認する (これはコミットに依存するので **逐次実行**)。

### 4. push (並列実行)

コミットが成功したら、続けて現在ブランチを `origin` に push する。

事前確認として次を並列実行して状態を把握する:

- `git rev-parse --abbrev-ref HEAD` — 現在のブランチ名
- `git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null || echo "no upstream"` — upstream の有無

push コマンドは状況に応じて切り替える:

- **upstream あり**: `git push`
- **upstream なし** (新規ブランチ): `git push -u origin <current-branch>`

その後 `git status` で push 後の状態 (ahead/behind 表示が消えていること) を確認する。

push 時の禁止事項:

- `--force` / `--force-with-lease` は **使わない** (ユーザが明示的に依頼した場合のみ)。
- `main` / `master` への force push は、たとえユーザが依頼しても警告してから実行する。
- push 先 (remote / ブランチ) を勝手に書き換えない。常に現在ブランチを同名で `origin` に push する。
- push が失敗した場合 (非 fast-forward 等) は安易に force しない。原因 (リモートに新しいコミットがある等) をユーザに伝えて指示を仰ぐ。

### 5. pre-commit / pre-push hook が失敗した場合

- `--amend` は **使わない**。hook が失敗したコミットは作られていないので、amend すると 1 つ前のコミットを書き換えてしまう。
- 原因を修正 → 再度 `git add` → **新しいコミット** を作る。
- `--no-verify` は使わない (ユーザが明示的に許可した場合を除く)。

## してはいけないこと

- `git config` は触らない。
- 変更がない (untracked も差分もない) 場合は空コミットを作らず、その旨を伝えて終了する (push もしない)。
- `git rebase -i` / `git add -i` など対話モードを必要とするコマンドは使わない。
- 別ブランチへのチェックアウトや、コミット対象外ファイルの破棄 (`git restore` / `git checkout --`) は行わない。
