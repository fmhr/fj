# fj

AtCoder Heuristic Contest (AHC) 向けのコマンドラインツール

## インストール

```bash
go install github.com/fmhr/fj@latest
```

## クイックスタート

```bash
# 1. テスターをダウンロード（AtCoderからZIPのURLを取得）
fj download https://img.atcoder.jp/ahcXXX/tools.zip

# 2. 初期化（設定ファイルとベストスコア追跡を作成）
fj init

# 3. テスト実行
fj test ./a.out
# または
fj test "go run main.go"
```

## コマンド

### `fj init`

設定ファイルとベストスコア追跡を初期化します。

```bash
# 基本的な初期化（fj/config.toml を生成）
fj init

# ベストスコア追跡を有効化（最大化問題）
fj init --best

# 最小化問題の場合
fj init --best --minimax=min
```

### `fj test` (エイリアス: `fj t`)

テストケースを実行します。

```bash
# 単一のテストケース実行（seed=0）
fj test ./a.out

# シード指定
fj test ./a.out -s 5

# 複数ケース実行
fj test ./a.out -n 100

# 並列実行
fj test ./a.out -n 100 -p 4

# ベストスコアを更新
fj test ./a.out -n 100 --best
```

| オプション | 短縮形 | 説明 |
|-----------|-------|------|
| `--seed` | `-s` | シード値（デフォルト: 0） |
| `--count` | `-n` | テストケース数（デフォルト: 1） |
| `--parallel` | `-p` | 並列実行数（デフォルト: 1） |

### `fj download` (エイリアス: `fj d`)

テスターのZIPファイルをダウンロードして展開します。

```bash
fj download https://img.atcoder.jp/ahcXXX/tools.zip
```

### `fj info`

現在の設定情報を表示します。

```bash
fj info
```

### `fj version` (エイリアス: `fj v`)

バージョン情報を表示します。

## グローバルオプション

| オプション | 説明 |
|-----------|------|
| `--debug` | デバッグモード（ログにファイル名:行番号を表示） |
| `--best` | ベストスコアを更新 |
| `--show-stderr` | 実行プログラムのstderr出力を表示 |
| `--cloud` | クラウドモードを有効化 |
| `--json` | JSON形式で出力 |
| `--csv <filename>` | CSV形式でファイルに出力 |

## 設定ファイル

`fj init` を実行すると `fj/config.toml` が生成されます。

```toml
Reactive = false
TesterPath = 'tools/target/release/tester'
VisPath = 'tools/target/release/vis'
GenPath = 'tools/target/release/gen'
InfilePath = 'tools/in/'
OutfilePath = 'out/'
TimeLimitMS = 5000
```

| 設定項目 | 説明 |
|---------|------|
| `Reactive` | リアクティブ問題かどうか |
| `TesterPath` | リアクティブ問題のテスターパス |
| `VisPath` | ビジュアライザ/スコア計算プログラムのパス |
| `InfilePath` | 入力ファイルのディレクトリ |
| `OutfilePath` | 出力ファイルのディレクトリ |
| `TimeLimitMS` | タイムアウト時間（ミリ秒） |

## ディレクトリ構成

```
project/
├── fj/
│   ├── config.toml      # 設定ファイル
│   └── best_score.json  # ベストスコア記録（--best使用時）
├── tools/
│   ├── in/              # 入力ファイル
│   └── target/release/
│       ├── tester       # リアクティブ用テスター
│       └── vis          # ビジュアライザ
├── out/                 # 出力ファイル
└── a.out                # 実行ファイル
```

## ライセンス

MIT License
